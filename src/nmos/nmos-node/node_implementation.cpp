#include <iostream>
#include <cstdlib>
#include <thread>
#include <unordered_map>
#include <algorithm>

#include "node_implementation.h"

#include <boost/range/adaptor/filtered.hpp>
#include <boost/range/adaptor/transformed.hpp>
#include <boost/range/algorithm/find.hpp>
#include <boost/range/algorithm/find_first_of.hpp>
#include <boost/range/algorithm/find_if.hpp>
#include <boost/range/algorithm_ext/push_back.hpp>
#include <boost/range/irange.hpp>
#include <boost/range/join.hpp>
#include "pplx/pplx_utils.h" // for pplx::complete_after, etc.
#include "cpprest/host_utils.h"
#ifdef HAVE_LLDP
#include "lldp/lldp_manager.h"
#endif
#include "nmos/activation_mode.h"
#include "nmos/capabilities.h"
#include "nmos/channels.h"
#include "nmos/channelmapping_resources.h"
#include "nmos/clock_name.h"
#include "nmos/colorspace.h"
#include "nmos/connection_resources.h"
#include "nmos/connection_events_activation.h"
#include "nmos/control_protocol_resources.h"
#include "nmos/control_protocol_resource.h"
#include "nmos/control_protocol_state.h"
#include "nmos/control_protocol_utils.h"
#include "nmos/events_resources.h"
#include "nmos/format.h"
#include "nmos/group_hint.h"
#include "nmos/interlace_mode.h"
#include "nmos/is12_versions.h"  // for IS-12 gain control
#ifdef HAVE_LLDP
#include "nmos/lldp_manager.h"
#endif
#include "nmos/media_type.h"
#include "nmos/model.h"
#include "nmos/node_interfaces.h"
#include "nmos/node_resource.h"
#include "nmos/node_resources.h"
#include "nmos/node_server.h"
#include "nmos/random.h"
#include "nmos/sdp_utils.h"
#include "nmos/slog.h"
#include "nmos/st2110_21_sender_type.h"
#include "nmos/system_resources.h"
#include "nmos/transfer_characteristic.h"
#include "nmos/transport.h"
#include "nmos/video_jxsv.h"
#include "sdp/sdp.h"

#include "FFmpeg_wrapper_client.h"

namespace tracker{
    std::unordered_map<nmos::id, Stream> resource_stream_map;
    Stream get_stream_info(const nmos::id& resource_id) {
        auto it = resource_stream_map.find(resource_id);
        if (it != resource_stream_map.end()) {
            return it->second;
        }
        throw std::runtime_error("Resource ID not found");
    }

    int get_stream_index(const nmos::id& resource_id) {
        auto it = resource_stream_map.find(resource_id);
        if (it != resource_stream_map.end()) {
            return std::distance(resource_stream_map.begin(), it);
        }
        throw std::runtime_error("Resource ID not found");
    }
    void add_stream_info(const nmos::id& resource_id, const Stream& stream_info) {
        resource_stream_map.emplace(resource_id, stream_info);
    }

    std::vector<Stream> get_file_streams_receivers(const Config& config) {
        std::vector<Stream> file_streams;
        if (config.receivers.empty()) {
            std::cout<< "No receivers in config in get_file_streams_receivers" << std::endl;
        }
        std::copy_if(config.receivers.begin(), config.receivers.end(), std::back_inserter(file_streams), [](const Stream& stream) {
            return stream.stream_type.type == stream_type::file;
        });
        return file_streams;
    }

    std::vector<Stream> get_file_streams_senders(const Config& config) {
        std::vector<Stream> file_streams;
        if (config.senders.empty()) {
            std::cout<< "No senders in config in get_file_streams_senders" << std::endl;
        }
        std::copy_if(config.senders.begin(), config.senders.end(), std::back_inserter(file_streams), [](const Stream& stream) {
            return stream.stream_type.type == stream_type::file;
        });
        return file_streams;
    }
}

namespace grpc {
// Function to execute the grpc client logic for ffmpeg
void sendDataToFfmpeg(const std::string& interface, const std::string& port, const Config &configParams) {
    CmdPassClient obj(interface, port);
    std::vector<std::pair<std::string,std::string>> commitedConfigs = commitConfigs(configParams);
    obj.FFmpegCmdExec(commitedConfigs);
}
} // namespace grpc

namespace impl
{
    // custom logging category for the example node implementation thread
    namespace categories
    {
        const nmos::category node_implementation{ "node_implementation" };
    }

    // custom settings for the example node implementation
    namespace fields
    {
        // node_tags, device_tags: used in resource tags fields
        // "Each tag has a single key, but MAY have multiple values."
        // See https://specs.amwa.tv/is-04/releases/v1.3.2/docs/APIs_-_Common_Keys.html#tags
        // {
        //     "tag_1": [ "tag_1_value_1", "tag_1_value_2" ],
        //     "tag_2": [ "tag_2_value_1" ]
        // }
        const web::json::field_as_value_or node_tags{ U("node_tags"), web::json::value::object() };
        const web::json::field_as_value_or device_tags{ U("device_tags"), web::json::value::object() };

        // activate_senders: controls whether to activate senders on start up (true, default) or not (false)
        const web::json::field_as_bool_or activate_senders{ U("activate_senders"), true };

        const web::json::field_as_array sender{ U("sender") };
        const web::json::field_as_array receiver{ U("receiver") };

        // sender_payload_type, receiver_payload_type: controls the payload_type of senders and receivers
        // TODO: change the reference to Config by stream sender or receiver
        const web::json::field_as_integer_or sender_payload_type{ U("sender_payload_type"), 112 };
        const web::json::field_as_integer_or receiver_payload_type{ U("receiver_payload_type"), 112 };

        // IP address and port to connect to ffmpeg grpc service in order to pass properties when another NMOS node
        // is connected to this NMOS node for BCS pipeline
        const web::json::field_as_string_or ffmpeg_grpc_server_address{ U("ffmpeg_grpc_server_address"), "localhost"};
        const web::json::field_as_string_or ffmpeg_grpc_server_port{ U("ffmpeg_grpc_server_port"), "50051"};

        // // default values for the node implementation
        const web::json::field_as_string_or default_source_ip{ U("default_source_ip"), "192.168.2.1"};
        const web::json::field_as_integer_or default_source_port{ U("default_source_port"), 5005};
        const web::json::field_as_integer_or default_destination_port{ U("default_destination_port"), 5000};
        const web::json::field_as_string_or default_destination_ip{ U("default_destination_ip"), "192.168.2.2"};
        const web::json::field_as_string_or default_interface_ip{ U("default_interface_ip"), "192.168.2.2"};

        const web::json::field_as_value_or frame_rate{ U("frame_rate"), web::json::value_of({
            { nmos::fields::numerator, 25 },
            { nmos::fields::denominator, 1 }
        }) };

         // frame_width, frame_height: control the frame_width and frame_height of video flows
        const web::json::field_as_integer_or frame_width{ U("frame_width"), 1920 };
        const web::json::field_as_integer_or frame_height{ U("frame_height"), 1080 };

        // In curret use case, below params are default values
        // interlace_mode: controls the interlace_mode of video flows, see nmos::interlace_mode
        // when omitted, a default of "progressive" or "interlaced_tff" is used based on the frame_rate, etc.
        const web::json::field_as_string interlace_mode{ U("interlace_mode") };

        // colorspace: controls the colorspace of video flows, see nmos::colorspace
        const web::json::field_as_string_or colorspace{ U("colorspace"), U("BT709") };

        // transfer_characteristic: controls the transfer characteristic system of video flows, see nmos::transfer_characteristic
        const web::json::field_as_string_or transfer_characteristic{ U("transfer_characteristic"), U("SDR") };

        // color_sampling: controls the color (sub-)sampling mode of video flows, see sdp::sampling
        const web::json::field_as_string_or color_sampling{ U("color_sampling"), U("YCbCr-4:2:2") };

        // component_depth: controls the bits per component sample of video flows
        const web::json::field_as_integer_or component_depth{ U("component_depth"), 10 };

        // video_type: media type of video flows, e.g. "video/raw" or "video/jxsv", see nmos::media_types
        const web::json::field_as_string_or video_type{ U("video_type"), U("video/raw") };

        // channel_count: controls the number of channels in audio sources
        const web::json::field_as_integer_or channel_count{ U("channel_count"), 4 };

        // smpte2022_7: controls whether senders and receivers have one leg (false) or two legs (true, default)
        const web::json::field_as_bool_or smpte2022_7{ U("smpte2022_7"), false };
    }

    nmos::interlace_mode get_interlace_mode(const nmos::rational& frame_rate, uint32_t frame_height, const nmos::settings& settings);

    // the different kinds of 'port' (standing for the format/media type/event type) implemented by the example node
    // each 'port' of the example node has a source, flow, sender and/or compatible receiver
    DEFINE_STRING_ENUM(port)
    namespace ports
    {
        // video/raw, video/jxsv, etc.
        const port video{ U("v") };
        // audio/L24
        const port audio{ U("a") };
        // video/smpte291
        const port data{ U("d") };
        // video/SMPTE2022-6
        const port mux{ U("m") };

        const std::vector<port> rtp{ video, audio, data, mux };
    }

    const std::vector<nmos::channel> channels_repeat{
        { U("Left Channel"), nmos::channel_symbols::L },
        { U("Right Channel"), nmos::channel_symbols::R },
        { U("Center Channel"), nmos::channel_symbols::C },
        { U("Low Frequency Effects Channel"), nmos::channel_symbols::LFE }
    };

    // find interface with the specified address
    std::vector<web::hosts::experimental::host_interface>::const_iterator find_interface(const std::vector<web::hosts::experimental::host_interface>& interfaces, const utility::string_t& address);

    // generate repeatable ids for the example node's resources
    nmos::id make_id(const nmos::id& seed_id, const nmos::type& type, const port& port = {}, int index = 0);
    std::vector<nmos::id> make_ids(const nmos::id& seed_id, const nmos::type& type, const port& port, int how_many = 1);
    std::vector<nmos::id> make_ids(const nmos::id& seed_id, const nmos::type& type, const std::vector<port>& ports, int how_many = 1);
    std::vector<nmos::id> make_ids(const nmos::id& seed_id, const std::vector<nmos::type>& types, const std::vector<port>& ports, int how_many = 1);

    // generate a repeatable source-specific multicast address for each leg of a sender
    utility::string_t make_source_specific_multicast_address_v4(const nmos::id& id, int leg = 0);

    // add a selection of parents to a source or flow
    void insert_parents(nmos::resource& resource, const nmos::id& seed_id, const port& port, int index);

    // add a helpful suffix to the label of a sub-resource for the example node
    void set_label_description(nmos::resource& resource, const port& port, int index);

    // add an example "natural grouping" hint to a sender or receiver
    void insert_group_hint(nmos::resource& resource, const port& port, int index);
}

// forward declarations for node_implementation_thread
void node_implementation_init(nmos::node_model& model, nmos::experimental::control_protocol_state& control_protocol_state, ConfigManager& config_manager, slog::base_gate& gate);
void node_implementation_run(nmos::node_model& model, slog::base_gate& gate);
nmos::connection_resource_auto_resolver make_node_implementation_auto_resolver(const nmos::settings& settings, ConfigManager& config_manager, slog::base_gate& gate);
nmos::connection_sender_transportfile_setter make_node_implementation_transportfile_setter(const nmos::resources& node_resources, const nmos::settings& settings, ConfigManager& config_manager, slog::base_gate& gate);

struct node_implementation_init_exception {};

// This is an example of how to integrate the nmos-cpp library with a device-specific underlying implementation.
// It constructs and inserts a node resource and some sub-resources into the model, based on the model settings,
// starts background tasks to emit regular events from the temperature event source, and then waits for shutdown.
void node_implementation_thread(nmos::node_model& model, nmos::experimental::control_protocol_state& control_protocol_state, ConfigManager& config_manager, slog::base_gate& gate_)
{
    nmos::details::omanip_gate gate{ gate_, nmos::stash_category(impl::categories::node_implementation) };

    try
    {
        node_implementation_init(model, control_protocol_state, config_manager, gate);
        node_implementation_run(model, gate);
    }
    catch (const node_implementation_init_exception&)
    {
        // node_implementation_init writes the log message
    }
    catch (const web::json::json_exception& e)
    {
        // most likely from incorrect value types in the command line settings
        slog::log<slog::severities::error>(gate, SLOG_FLF) << "JSON error: " << e.what();
    }
    catch (const std::system_error& e)
    {
        slog::log<slog::severities::error>(gate, SLOG_FLF) << "System error: " << e.what() << " [" << e.code() << "]";
    }
    catch (const std::runtime_error& e)
    {
        slog::log<slog::severities::error>(gate, SLOG_FLF) << "Implementation error: " << e.what();
    }
    catch (const std::exception& e)
    {
        slog::log<slog::severities::error>(gate, SLOG_FLF) << "Unexpected exception: " << e.what();
    }
    catch (...)
    {
        slog::log<slog::severities::severe>(gate, SLOG_FLF) << "Unexpected unknown exception";
    }
}

void node_implementation_init(nmos::node_model& model, nmos::experimental::control_protocol_state& control_protocol_state, ConfigManager& config_manager, slog::base_gate& gate)
{
    using web::json::value;
    using web::json::value_from_elements;
    using web::json::value_of;

    auto lock = model.write_lock(); // in order to update the resources

    const auto seed_id = nmos::experimental::fields::seed_id(model.settings);
    const auto node_id = impl::make_id(seed_id, nmos::types::node);
    const auto device_id = impl::make_id(seed_id, nmos::types::device);

    const auto ffmpeg_grpc_server_address = impl::fields::ffmpeg_grpc_server_address(model.settings);
    const auto ffmpeg_grpc_server_port = impl::fields::ffmpeg_grpc_server_port(model.settings);

    const auto default_source_ip = impl::fields::default_source_ip(model.settings);
    const auto default_source_port = impl::fields::default_source_port(model.settings);
    const auto default_destination_ip = impl::fields::default_destination_ip(model.settings);
    const auto default_destination_port = impl::fields::default_destination_port(model.settings);
    const auto default_interface_ip = impl::fields::default_interface_ip(model.settings);

    //seder payload is right know global for all senders
    const auto sender_payload_type = impl::fields::sender_payload_type(model.settings);

    auto configIntel = config_manager.get_config();
    auto sender_arr_length = configIntel.senders.size();
    auto sender_arr = configIntel.senders;
    auto receiver_arr_length = configIntel.receivers.size();
    auto receiver_arr = configIntel.receivers;

    const std::vector<impl::port> media_ports = { impl::ports::video };
 
    //generic values for whole node
    //const auto interlace_mode = impl::get_interlace_mode(model.settings);
    const auto frame_rate = nmos::parse_rational(impl::fields::frame_rate(model.settings));
    const auto colorspace = nmos::colorspace{ impl::fields::colorspace(model.settings) };
    const auto transfer_characteristic = nmos::transfer_characteristic{ impl::fields::transfer_characteristic(model.settings) };
    const auto sampling = sdp::sampling{ impl::fields::color_sampling(model.settings) };
    const auto bit_depth = impl::fields::component_depth(model.settings);
    // const auto video_type = nmos::media_type{ impl::fields::video_type(model.settings) };
    const auto channel_count = impl::fields::channel_count(model.settings);
    const auto smpte2022_7 = impl::fields::smpte2022_7(model.settings);

    // for now, some typical values for video/jxsv, based on VSF TR-08:2022
    // see https://vsf.tv/download/technical_recommendations/VSF_TR-08_2022-04-20.pdf
    const auto profile = nmos::profiles::High444_12;
    // const auto level = nmos::get_video_jxsv_level(frame_rate, frame_width, frame_height);
    const auto sublevel = nmos::sublevels::Sublev3bpp;
    const auto max_bits_per_pixel = 4.0; // min coding efficiency
    const auto bits_per_pixel = 2.0;
    const auto transport_bit_rate_factor = 1.05;

    // any delay between updates to the model resources is unnecessary unless for debugging purposes
    const unsigned int delay_millis{ 0 };

    // it is important that the model be locked before inserting, updating or deleting a resource
    // and that the the node behaviour thread be notified after doing so
    const auto insert_resource_after = [&model, &lock](unsigned int milliseconds, nmos::resources& resources, nmos::resource&& resource, slog::base_gate& gate)
    {
        if (nmos::details::wait_for(model.shutdown_condition, lock, std::chrono::milliseconds(milliseconds), [&] { return model.shutdown; })) return false;

        const std::pair<nmos::id, nmos::type> id_type{ resource.id, resource.type };
        const bool success = insert_resource(resources, std::move(resource)).second;

        if (success)
            slog::log<slog::severities::info>(gate, SLOG_FLF) << "Updated model with " << id_type;
        else
            slog::log<slog::severities::severe>(gate, SLOG_FLF) << "Model update error: " << id_type;

        slog::log<slog::severities::too_much_info>(gate, SLOG_FLF) << "Notifying node behaviour thread";
        model.notify();

        return success;
    };

    const auto resolve_auto = make_node_implementation_auto_resolver(model.settings, config_manager, gate);
    const auto set_transportfile = make_node_implementation_transportfile_setter(model.node_resources, model.settings, config_manager, gate);

    const auto clocks = web::json::value_of({ nmos::make_internal_clock(nmos::clock_names::clk0) });
    // filter network interfaces to those that correspond to the specified host_addresses
    const auto host_interfaces = nmos::get_host_interfaces(model.settings);
    const auto interfaces = nmos::experimental::node_interfaces(host_interfaces);

    {
        auto node = nmos::make_node(node_id, clocks, nmos::make_node_interfaces(interfaces), model.settings);
        node.data[nmos::fields::tags] = impl::fields::node_tags(model.settings);
        if (!insert_resource_after(delay_millis, model.node_resources, std::move(node), gate)) throw node_implementation_init_exception();
    }

#ifdef HAVE_LLDP
    // LLDP manager for advertising server identity, capabilities, and discovering neighbours on a local area network
    slog::log<slog::severities::info>(gate, SLOG_FLF) << "Attempting to configure LLDP";
    auto lldp_manager = nmos::experimental::make_lldp_manager(model, interfaces, true, gate);
    // hm, open may potentially throw?
    lldp::lldp_manager_guard lldp_manager_guard(lldp_manager);
#endif

    // prepare interface bindings for all senders and receivers
    const auto& host_address = nmos::fields::host_address(model.settings);
    // the interface corresponding to the host address is used for the example node's WebSocket senders and receivers
    const auto host_interface_ = impl::find_interface(host_interfaces, host_address);
    if (host_interfaces.end() == host_interface_)
    {
        slog::log<slog::severities::severe>(gate, SLOG_FLF) << "No network interface corresponding to host_address?";
        throw node_implementation_init_exception();
    }
    // const auto& host_interface = *host_interface_;
    // hmm, should probably add a custom setting to control the primary and secondary interfaces for the example node's RTP senders and receivers
    // rather than just picking the one(s) corresponding to the first and last of the specified host addresses
    const auto& primary_address = model.settings.has_field(nmos::fields::host_addresses) ? web::json::front(nmos::fields::host_addresses(model.settings)).as_string() : host_address;
    const auto& secondary_address = model.settings.has_field(nmos::fields::host_addresses) ? web::json::back(nmos::fields::host_addresses(model.settings)).as_string() : host_address;
    const auto primary_interface_ = impl::find_interface(host_interfaces, primary_address);
    const auto secondary_interface_ = impl::find_interface(host_interfaces, secondary_address);
    if (host_interfaces.end() == primary_interface_ || host_interfaces.end() == secondary_interface_)
    {
        slog::log<slog::severities::severe>(gate, SLOG_FLF) << "No network interface corresponding to one of the host_addresses?";
        throw node_implementation_init_exception();
    }
    const auto& primary_interface = *primary_interface_;
    const auto& secondary_interface = *secondary_interface_;
    const auto interface_names = smpte2022_7
        ? std::vector<utility::string_t>{ primary_interface.name, secondary_interface.name }
        : std::vector<utility::string_t>{ primary_interface.name };

    {
        // For simplified NMOS and BCS needs, only one device = pipeline is required.
        slog::log<slog::severities::info>(gate, SLOG_FLF) << "DEVICE";
        slog::log<slog::severities::info>(gate, SLOG_FLF) << "SENDERS_TOTAL = " << sender_arr_length;
        slog::log<slog::severities::info>(gate, SLOG_FLF) << "RECEIVERS_TOTAL = " << receiver_arr_length;

        auto sender_ids = impl::make_ids(seed_id, nmos::types::sender, media_ports, sender_arr_length);
        auto receiver_ids = impl::make_ids(seed_id, nmos::types::receiver, media_ports, receiver_arr_length);

        auto device = nmos::make_device(device_id, node_id, sender_ids, receiver_ids, model.settings);

        device.data[nmos::fields::tags] = impl::fields::device_tags(model.settings);

        if (!insert_resource_after(delay_millis, model.node_resources, std::move(device), gate)) throw node_implementation_init_exception();
    }

    // currently only media_port video is supported, in next iterations, the support of audio will be implemented
    for (const auto& port : media_ports)
    {
        for (int index = 0; index < sender_arr_length; ++index)
        {
            if (sender_arr[index].stream_type.type == stream_type::file)
            {
                std::cout<< "Sender stream type is file" << std::endl;
                continue;
            }
            const auto source_id = impl::make_id(seed_id, nmos::types::source, port, index);
            const auto flow_id = impl::make_id(seed_id, nmos::types::flow, port, index);
            const auto sender_id = impl::make_id(seed_id, nmos::types::sender, port, index);

            auto senderDefinition = sender_arr[index];

            const auto frame_rate_json_format = web::json::value_of({
                { nmos::fields::numerator, config_manager.get_framerate(senderDefinition).first },
                { nmos::fields::denominator, config_manager.get_framerate(senderDefinition).second }
            });

            const auto frame_rate_parsed_rational = nmos::parse_rational(frame_rate_json_format);
            const auto frame_w = senderDefinition.payload.video.frame_width;
            const auto frame_h = senderDefinition.payload.video.frame_height;
            const auto level = nmos::get_video_jxsv_level(frame_rate_parsed_rational, frame_w, frame_h);
            const auto tx_interlace_mode = impl::get_interlace_mode(frame_rate_parsed_rational, frame_h, model.settings);

            nmos::media_type video_type;
            if (senderDefinition.payload.video.video_type == "rawvideo") {
                video_type = nmos::media_types::video_raw;
            } else if (senderDefinition.payload.video.video_type == "jxsv") {
                video_type = nmos::media_types::video_jxsv;
            } else {
                // https://specs.amwa.tv/is-04/releases/v1.2.0/APIs/schemas/with-refs/flow_video_coded.html
                video_type = nmos::media_type{utility::s2us(senderDefinition.payload.video.video_type)};
            }

            nmos::resource source;
            if (impl::ports::video == port)
            {
                source = nmos::make_video_source(source_id, device_id, nmos::clock_names::clk0, frame_rate_parsed_rational, model.settings);
            }
            else if (impl::ports::audio == port) // not yet supported
            {
                const auto channels = boost::copy_range<std::vector<nmos::channel>>(boost::irange(0, channel_count) | boost::adaptors::transformed([&](const int& index)
                {
                    return impl::channels_repeat[index % (int)impl::channels_repeat.size()];
                }));

                source = nmos::make_audio_source(source_id, device_id, nmos::clock_names::clk0, frame_rate_parsed_rational, channels, model.settings);
            }
            else if (impl::ports::data == port) // not yet supported
            {
                source = nmos::make_data_source(source_id, device_id, nmos::clock_names::clk0, frame_rate_parsed_rational, model.settings);
            }
            else if (impl::ports::mux == port) // not yet supported
            {
                source = nmos::make_mux_source(source_id, device_id, nmos::clock_names::clk0, frame_rate_parsed_rational, model.settings);
            }
            impl::insert_parents(source, seed_id, port, index);
            impl::set_label_description(source, port, index);

            nmos::resource flow;
            if (impl::ports::video == port)
            {
                if (nmos::media_types::video_raw == video_type)
                {
                    flow = nmos::make_raw_video_flow(
                        flow_id, source_id, device_id,
                        frame_rate_parsed_rational,
                        frame_w, frame_h, tx_interlace_mode,
                        colorspace, transfer_characteristic, sampling, bit_depth,
                        model.settings
                    );
                }
                else if (nmos::media_types::video_jxsv == video_type)
                {
                    flow = nmos::make_video_jxsv_flow(
                        flow_id, source_id, device_id,
                        frame_rate_parsed_rational,
                        frame_w, frame_h, tx_interlace_mode,
                        colorspace, transfer_characteristic, sampling, bit_depth,
                        profile, level, sublevel, bits_per_pixel,
                        model.settings
                    );
                }
                else
                {
                    flow = nmos::make_coded_video_flow(
                        flow_id, source_id, device_id,
                        frame_rate_parsed_rational,
                        frame_w, frame_h, tx_interlace_mode,
                        colorspace, transfer_characteristic, sampling, bit_depth,
                        video_type,
                        model.settings
                    );
                }
            }
            else if (impl::ports::audio == port)
            {
                flow = nmos::make_raw_audio_flow(flow_id, source_id, device_id, 48000, 24, model.settings);
                // add optional grain_rate
                flow.data[nmos::fields::grain_rate] = nmos::make_rational(frame_rate_parsed_rational);
            }
            else if (impl::ports::data == port)
            {
                nmos::did_sdid timecode{ 0x60, 0x60 };
                flow = nmos::make_sdianc_data_flow(flow_id, source_id, device_id, { timecode }, model.settings);
                // add optional grain_rate
                flow.data[nmos::fields::grain_rate] = nmos::make_rational(frame_rate_parsed_rational);
            }
            else if (impl::ports::mux == port)
            {
                flow = nmos::make_mux_flow(flow_id, source_id, device_id, model.settings);
                // add optional grain_rate
                flow.data[nmos::fields::grain_rate] = nmos::make_rational(frame_rate_parsed_rational);
            }
            impl::insert_parents(flow, seed_id, port, index);
            impl::set_label_description(flow, port, index);

            // set_transportfile needs to find the matching source and flow for the sender, so insert these first
            if (!insert_resource_after(delay_millis, model.node_resources, std::move(source), gate)) throw node_implementation_init_exception();
            if (!insert_resource_after(delay_millis, model.node_resources, std::move(flow), gate)) throw node_implementation_init_exception();

            const auto manifest_href = nmos::experimental::make_manifest_api_manifest(sender_id, model.settings);
            std::cout<< "Make sender"<< std::endl;
            auto sender = nmos::make_sender(sender_id, flow_id, nmos::transports::rtp, device_id, manifest_href.to_string(), interface_names, model.settings);
            tracker::add_stream_info(sender_id,senderDefinition);
            // hm, could add nmos::make_video_jxsv_sender to encapsulate this?
            if (impl::ports::video == port && nmos::media_types::video_jxsv == video_type)
            {
                // additional attributes required by BCP-006-01
                // see https://specs.amwa.tv/bcp-006-01/branches/v1.0-dev/docs/NMOS_With_JPEG_XS.html#senders
                const auto format_bit_rate = nmos::get_video_jxsv_bit_rate(frame_rate_parsed_rational, frame_w, frame_h, bits_per_pixel);
                // round to nearest Megabit/second per examples in VSF TR-08:2022
                const auto transport_bit_rate = uint64_t(transport_bit_rate_factor * format_bit_rate / 1e3 + 0.5) * 1000;
                sender.data[nmos::fields::bit_rate] = value(transport_bit_rate);
                sender.data[nmos::fields::st2110_21_sender_type] = value(nmos::st2110_21_sender_types::type_N.name);
            }
            impl::set_label_description(sender, port, index);
            impl::insert_group_hint(sender, port, index);

            auto connection_sender = nmos::make_connection_rtp_sender(sender_id, smpte2022_7);

            auto sender_stream_number = tracker::get_stream_index(sender_id);

            connection_sender.data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::source_ip] = web::json::value::string(impl::fields::default_source_ip(model.settings));
            connection_sender.data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::destination_ip] = web::json::value::string(impl::fields::default_destination_ip(model.settings));
            connection_sender.data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::source_port] = web::json::value::number(impl::fields::default_source_port(model.settings) + sender_stream_number);
            connection_sender.data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::destination_port] = web::json::value::number(impl::fields::default_destination_port(model.settings)+sender_stream_number);

            if (impl::fields::activate_senders(model.settings))
            {
                // initialize this sender with a scheduled activation, e.g. to enable the IS-05-01 test suite to run immediately
                auto& staged = connection_sender.data[nmos::fields::endpoint_staged];
                staged[nmos::fields::master_enable] = value::boolean(true);
                staged[nmos::fields::activation] = value_of({
                    { nmos::fields::mode, nmos::activation_modes::activate_scheduled_relative.name },
                    { nmos::fields::requested_time, U("0:0") },
                    { nmos::fields::activation_time, nmos::make_version() }
                });
            }

            if (!insert_resource_after(delay_millis, model.node_resources, std::move(sender), gate)) throw node_implementation_init_exception();
            if (!insert_resource_after(delay_millis, model.connection_resources, std::move(connection_sender), gate)) throw node_implementation_init_exception();
        }
    }

    for (const auto& port : media_ports) {
        for (int index = 0; index < receiver_arr_length; ++index){
            if (receiver_arr[index].stream_type.type == stream_type::file)
            {
                std::cout<< "Receiver stream type is file" << std::endl;
                continue;
            }

            const auto receiver_id = impl::make_id(seed_id, nmos::types::receiver, port, index);
            auto receiverDefinition = receiver_arr[index];
            const auto rx_frame_rate_json_format = web::json::value_of({
                { nmos::fields::numerator, config_manager.get_framerate(receiverDefinition).first },
                { nmos::fields::denominator, config_manager.get_framerate(receiverDefinition).second }
            });
            const auto rx_frame_rate_parsed_rational = nmos::parse_rational(rx_frame_rate_json_format);
            const auto frame_w_r = receiverDefinition.payload.video.frame_width;
            const auto frame_h_r = receiverDefinition.payload.video.frame_height;
            const auto level = nmos::get_video_jxsv_level(rx_frame_rate_parsed_rational, frame_w_r, frame_h_r);
            const auto rx_interlace_mode = impl::get_interlace_mode(rx_frame_rate_parsed_rational, frame_h_r, model.settings);
            nmos::media_type video_type;
            if (receiverDefinition.payload.video.video_type == "rawvideo") {
                video_type = nmos::media_types::video_raw;
            } else if (receiverDefinition.payload.video.video_type == "jxsv") {
                video_type = nmos::media_types::video_jxsv;
            } else {
                // https://specs.amwa.tv/is-04/releases/v1.2.0/APIs/schemas/with-refs/flow_video_coded.html
                video_type = nmos::media_type{utility::s2us(receiverDefinition.payload.video.video_type)};
            }
            nmos::resource receiver;
            if (impl::ports::video == port)
            {
                receiver = nmos::make_receiver(receiver_id, device_id, nmos::transports::rtp, interface_names, nmos::formats::video, { video_type }, model.settings);
                tracker::add_stream_info(receiver_id, receiverDefinition);
                // add an example constraint set; these should be completed fully!
                if (nmos::media_types::video_raw == video_type)
                {
                    const auto interlace_modes = nmos::interlace_modes::progressive != rx_interlace_mode
                        ? std::vector<utility::string_t>{ nmos::interlace_modes::interlaced_bff.name, nmos::interlace_modes::interlaced_tff.name, nmos::interlace_modes::interlaced_psf.name }
                        : std::vector<utility::string_t>{ nmos::interlace_modes::progressive.name };
                    receiver.data[nmos::fields::caps][nmos::fields::constraint_sets] = value_of({
                        value_of({
                            { nmos::caps::format::grain_rate, nmos::make_caps_rational_constraint({ rx_frame_rate_parsed_rational }) },
                            { nmos::caps::format::frame_width, nmos::make_caps_integer_constraint({ frame_w_r }) },
                            { nmos::caps::format::frame_height, nmos::make_caps_integer_constraint({ frame_h_r }) },
                            { nmos::caps::format::interlace_mode, nmos::make_caps_string_constraint(interlace_modes) },
                            { nmos::caps::format::color_sampling, nmos::make_caps_string_constraint({ "YCbCr-4:2:2" }) }
                        })
                    });
                }
                else if (nmos::media_types::video_jxsv == video_type)
                {
                    // some of the parameter constraints recommended by BCP-006-01
                    // see https://specs.amwa.tv/bcp-006-01/branches/v1.0-dev/docs/NMOS_With_JPEG_XS.html#receivers
                    const auto max_format_bit_rate = nmos::get_video_jxsv_bit_rate(rx_frame_rate_parsed_rational, frame_w_r, frame_h_r, max_bits_per_pixel);
                    // round to nearest Megabit/second per examples in VSF TR-08:2022
                    const auto max_transport_bit_rate = uint64_t(transport_bit_rate_factor * max_format_bit_rate / 1e3 + 0.5) * 1000;

                    receiver.data[nmos::fields::caps][nmos::fields::constraint_sets] = value_of({
                        value_of({
                            { nmos::caps::format::profile, nmos::make_caps_string_constraint({ profile.name }) },
                            { nmos::caps::format::level, nmos::make_caps_string_constraint({ level.name }) },
                            { nmos::caps::format::sublevel, nmos::make_caps_string_constraint({ nmos::sublevels::Sublev3bpp.name, nmos::sublevels::Sublev4bpp.name }) },
                            { nmos::caps::format::bit_rate, nmos::make_caps_integer_constraint({}, nmos::no_minimum<int64_t>(), (int64_t)max_format_bit_rate) },
                            { nmos::caps::transport::bit_rate, nmos::make_caps_integer_constraint({}, nmos::no_minimum<int64_t>(), (int64_t)max_transport_bit_rate) },
                            { nmos::caps::transport::packet_transmission_mode, nmos::make_caps_string_constraint({ nmos::packet_transmission_modes::codestream.name }) }
                        })
                    });
                }
                
                receiver.data[nmos::fields::version] = receiver.data[nmos::fields::caps][nmos::fields::version] = value(nmos::make_version());
            }
            else if (impl::ports::audio == port)
            {
                receiver = nmos::make_audio_receiver(receiver_id, device_id, nmos::transports::rtp, interface_names, 24, model.settings);
                // add some example constraint sets; these should be completed fully!
                receiver.data[nmos::fields::caps][nmos::fields::constraint_sets] = value_of({
                    value_of({
                        { nmos::caps::format::channel_count, nmos::make_caps_integer_constraint({}, 1, channel_count) },
                        { nmos::caps::format::sample_rate, nmos::make_caps_rational_constraint({ { 48000, 1 } }) },
                        { nmos::caps::format::sample_depth, nmos::make_caps_integer_constraint({ 16, 24 }) },
                        { nmos::caps::transport::packet_time, nmos::make_caps_number_constraint({ 0.125 }) }
                    }),
                    value_of({
                        { nmos::caps::meta::preference, -1 },
                        { nmos::caps::format::channel_count, nmos::make_caps_integer_constraint({}, 1, (std::min)(8, channel_count)) },
                        { nmos::caps::format::sample_rate, nmos::make_caps_rational_constraint({ { 48000, 1 } }) },
                        { nmos::caps::format::sample_depth, nmos::make_caps_integer_constraint({ 16, 24 }) },
                        { nmos::caps::transport::packet_time, nmos::make_caps_number_constraint({ 1 }) }
                    })
                });
                receiver.data[nmos::fields::version] = receiver.data[nmos::fields::caps][nmos::fields::version] = value(nmos::make_version());
            }
            else if (impl::ports::data == port)
            {
                receiver = nmos::make_sdianc_data_receiver(receiver_id, device_id, nmos::transports::rtp, interface_names, model.settings);
                // add an example constraint set; these should be completed fully!
                receiver.data[nmos::fields::caps][nmos::fields::constraint_sets] = value_of({
                    value_of({
                        { nmos::caps::format::grain_rate, nmos::make_caps_rational_constraint({ rx_frame_rate_parsed_rational }) }
                    })
                });
                receiver.data[nmos::fields::version] = receiver.data[nmos::fields::caps][nmos::fields::version] = value(nmos::make_version());
            }
            else if (impl::ports::mux == port)
            {
                receiver = nmos::make_mux_receiver(receiver_id, device_id, nmos::transports::rtp, interface_names, model.settings);
                // add an example constraint set; these should be completed fully!
                receiver.data[nmos::fields::caps][nmos::fields::constraint_sets] = value_of({
                    value_of({
                        { nmos::caps::format::grain_rate, nmos::make_caps_rational_constraint({ rx_frame_rate_parsed_rational }) }
                    })
                });
                receiver.data[nmos::fields::version] = receiver.data[nmos::fields::caps][nmos::fields::version] = value(nmos::make_version());
            }
            impl::set_label_description(receiver, port, index);
            impl::insert_group_hint(receiver, port, index);

            auto connection_receiver = nmos::make_connection_rtp_receiver(receiver_id, smpte2022_7);

            connection_receiver.data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::interface_ip] = web::json::value::string(impl::fields::default_interface_ip(model.settings));

            resolve_auto(receiver, connection_receiver, connection_receiver.data[nmos::fields::endpoint_active][nmos::fields::transport_params]);

            if (!insert_resource_after(delay_millis, model.node_resources, std::move(receiver), gate)) throw node_implementation_init_exception();
            if (!insert_resource_after(delay_millis, model.connection_resources, std::move(connection_receiver), gate)) throw node_implementation_init_exception();
        }
    }
}

void node_implementation_run(nmos::node_model& model, slog::base_gate& gate)
{
}

// Example System API node behaviour callback to perform application-specific operations when the global configuration resource changes
nmos::system_global_handler make_node_implementation_system_global_handler(nmos::node_model& model, slog::base_gate& gate)
{
    // this example uses the callback to update the settings
    // (an 'empty' std::function disables System API node behaviour)
    return [&](const web::uri& system_uri, const web::json::value& system_global)
    {
        if (!system_uri.is_empty())
        {
            slog::log<slog::severities::info>(gate, SLOG_FLF) << nmos::stash_category(impl::categories::node_implementation) << "New system global configuration discovered from the System API at: " << system_uri.to_string();

            // although this example immediately updates the settings, the effect is not propagated
            // in either Registration API behaviour or the senders' /transportfile endpoints until
            // an update to these is forced by other circumstances

            auto system_global_settings = nmos::parse_system_global_data(system_global).second;
            web::json::merge_patch(model.settings, system_global_settings, true);
        }
        else
        {
            slog::log<slog::severities::warning>(gate, SLOG_FLF) << nmos::stash_category(impl::categories::node_implementation) << "System global configuration is not discoverable";
        }
    };
}

// Example Registration API node behaviour callback to perform application-specific operations when the current Registration API changes
nmos::registration_handler make_node_implementation_registration_handler(slog::base_gate& gate)
{
    return [&](const web::uri& registration_uri)
    {
        if (!registration_uri.is_empty())
        {
            slog::log<slog::severities::info>(gate, SLOG_FLF) << nmos::stash_category(impl::categories::node_implementation) << "Started registered operation with Registration API at: " << registration_uri.to_string();
        }
        else
        {
            slog::log<slog::severities::warning>(gate, SLOG_FLF) << nmos::stash_category(impl::categories::node_implementation) << "Stopped registered operation";
        }
    };
}

// Example Connection API callback to parse "transport_file" during a PATCH /staged request
nmos::transport_file_parser make_node_implementation_transport_file_parser(slog::base_gate& gate)
{
    // this example uses a custom transport file parser to handle video/jxsv in addition to the core media types
    // otherwise, it could simply return &nmos::parse_rtp_transport_file
    // (if this callback is specified, an 'empty' std::function is not allowed)
    return [](const nmos::resource& receiver, const nmos::resource& connection_receiver, const utility::string_t& transport_file_type, const utility::string_t& transport_file_data, slog::base_gate& gate)
    {
        const auto validate_sdp_parameters = [&gate](const web::json::value& receiver, const nmos::sdp_parameters& sdp_params)
        {
            if (nmos::media_types::video_jxsv == nmos::get_media_type(sdp_params))
            {
                nmos::validate_video_jxsv_sdp_parameters(receiver, sdp_params);
            }
            else
            {
                // validate core media types, i.e., "video/raw", "audio/L", "video/smpte291" and "video/SMPTE2022-6"
                nmos::validate_sdp_parameters(receiver, sdp_params);
            }
        };
        return nmos::details::parse_rtp_transport_file(validate_sdp_parameters, receiver, connection_receiver, transport_file_type, transport_file_data, gate);
    };
}

// Example Connection API callback to perform application-specific validation of the merged /staged endpoint during a PATCH /staged request
nmos::details::connection_resource_patch_validator make_node_implementation_patch_validator(slog::base_gate& gate)
{
    // this example uses an 'empty' std::function because it does not need to do any validation
    // beyond what is expressed by the schemas and /constraints endpoint
    return{};
}


void filter_auto_default_transport_params_sender(web::json::value& transport_params, const nmos::settings& settings, const nmos::type& type_nmos, const web::json::value& data, nmos::id stream_id)
{
    auto sender_stream_number = tracker::get_stream_index(stream_id);
    auto is_auto = [](const web::json::value& v) {
        return v.is_string() && v.as_string() == U("auto");
    };

    web::json::value staged = data.has_field(nmos::fields::endpoint_staged) && data.at(nmos::fields::endpoint_staged).has_field(nmos::fields::transport_params)
        ? data.at(nmos::fields::endpoint_staged).at(nmos::fields::transport_params).at(0)
        : web::json::value();
    web::json::value active = data.has_field(nmos::fields::endpoint_active) && data.at(nmos::fields::endpoint_active).has_field(nmos::fields::transport_params)
        ? data.at(nmos::fields::endpoint_active).at(nmos::fields::transport_params).at(0)
        : web::json::value();

    if ((staged.has_field(nmos::fields::source_ip) && is_auto(staged[nmos::fields::source_ip])) ||
        (active.has_field(nmos::fields::source_ip) && is_auto(active[nmos::fields::source_ip])))
    {
        nmos::details::resolve_auto(transport_params[0], nmos::fields::source_ip, [&] {
            return web::json::value::string(impl::fields::default_source_ip(settings));
        });
    }
    else if (staged.has_field(nmos::fields::source_ip))
    {
        nmos::details::resolve_auto(transport_params[0], nmos::fields::source_ip, [&] {
            return staged[nmos::fields::source_ip];
        });
    }
    if ((staged.has_field(nmos::fields::destination_ip) && is_auto(staged[nmos::fields::destination_ip])) ||
        (active.has_field(nmos::fields::destination_ip) && is_auto(active[nmos::fields::destination_ip])))
    {
        nmos::details::resolve_auto(transport_params[0], nmos::fields::destination_ip, [&] {
            return web::json::value::string(impl::fields::default_destination_ip(settings));
        });
    } else if (staged.has_field(nmos::fields::destination_ip))
    {
        nmos::details::resolve_auto(transport_params[0], nmos::fields::destination_ip, [&] {
            return staged[nmos::fields::destination_ip];
        });
    }
    if ((staged.has_field(nmos::fields::source_port) && is_auto(staged[nmos::fields::source_port])) ||
        (active.has_field(nmos::fields::source_port) && is_auto(active[nmos::fields::source_port])))
    {
        nmos::details::resolve_auto(transport_params[0], nmos::fields::source_port, [&] {
            return web::json::value::number(impl::fields::default_source_port(settings)+sender_stream_number);
        });
    } else if (staged.has_field(nmos::fields::source_port))
    {
        nmos::details::resolve_auto(transport_params[0], nmos::fields::source_port, [&] {
            return staged[nmos::fields::source_port];
        });
    }
    if ((staged.has_field(nmos::fields::destination_port) && is_auto(staged[nmos::fields::destination_port])) ||
        (active.has_field(nmos::fields::destination_port) && is_auto(active[nmos::fields::destination_port])))
    {
        nmos::details::resolve_auto(transport_params[0], nmos::fields::destination_port, [&] {
            return web::json::value::number(impl::fields::default_destination_port(settings)+sender_stream_number);
        });
    } else if (staged.has_field(nmos::fields::destination_port))
    {
        nmos::details::resolve_auto(transport_params[0], nmos::fields::destination_port, [&] {
            return staged[nmos::fields::destination_port];
        });
    }
    nmos::resolve_rtp_auto(type_nmos, transport_params);
}
void filter_auto_default_transport_params_receiver(web::json::value& transport_params, const nmos::settings & settings, const nmos::type & type_nmos, const web::json::value& data)
{
    auto is_auto = [](const web::json::value& v) {
        return v.is_string() && v.as_string() == U("auto");
    };

    web::json::value staged = data.has_field(nmos::fields::endpoint_staged) && data.at(nmos::fields::endpoint_staged).has_field(nmos::fields::transport_params)
        ? data.at(nmos::fields::endpoint_staged).at(nmos::fields::transport_params).at(0)
        : web::json::value();
    web::json::value active = data.has_field(nmos::fields::endpoint_active) && data.at(nmos::fields::endpoint_active).has_field(nmos::fields::transport_params)
        ? data.at(nmos::fields::endpoint_active).at(nmos::fields::transport_params).at(0)
        : web::json::value();

    if ((staged.has_field(nmos::fields::interface_ip) && is_auto(staged[nmos::fields::interface_ip])) ||
        (active.has_field(nmos::fields::interface_ip) && is_auto(active[nmos::fields::interface_ip])))
    {
        nmos::details::resolve_auto(transport_params[0], nmos::fields::interface_ip, [&] {
            return web::json::value::string(impl::fields::default_interface_ip(settings));
        });
    }
    else if (staged.has_field(nmos::fields::interface_ip))
    {
        nmos::details::resolve_auto(transport_params[0], nmos::fields::interface_ip, [&] {
            return staged[nmos::fields::interface_ip];
        });
    }
nmos::resolve_rtp_auto(type_nmos, transport_params);
}

// Connection API activation callback to resolve "auto" values when /staged is transitioned to /active
nmos::connection_resource_auto_resolver make_node_implementation_auto_resolver(const nmos::settings& settings, ConfigManager& config_manager, slog::base_gate& gate)
{
    using web::json::value;

    const auto seed_id = nmos::experimental::fields::seed_id(settings);
    const auto device_id = impl::make_id(seed_id, nmos::types::device);

    auto configIntel = config_manager.get_config();
    auto sender_arr_length = configIntel.senders.size();
    auto sender_arr = configIntel.senders;
    auto receiver_arr_length = configIntel.receivers.size();
    auto receiver_arr = configIntel.receivers;

    const std::vector<impl::port> media_ports = { impl::ports::video };

    const auto rtp_sender_ids = impl::make_ids(seed_id, nmos::types::sender, media_ports, sender_arr_length);
    const auto rtp_receiver_ids = impl::make_ids(seed_id, nmos::types::receiver, media_ports, receiver_arr_length);

    return [rtp_sender_ids, rtp_receiver_ids, &gate, &settings](const nmos::resource& resource, const nmos::resource& connection_resource, value& transport_params)
    {
        const std::pair<nmos::id, nmos::type> id_type{ connection_resource.id, connection_resource.type };

        if (rtp_sender_ids.end() != boost::range::find(rtp_sender_ids, id_type.first))
        {
            // senders
            auto data=connection_resource.data;
            filter_auto_default_transport_params_sender(transport_params, settings, id_type.second, data, id_type.first);
        }
        else if (rtp_receiver_ids.end() != boost::range::find(rtp_receiver_ids, id_type.first))
        {
            // receivers
            auto data=connection_resource.data;
            filter_auto_default_transport_params_receiver(transport_params, settings, id_type.second, data);
        }
    };
}

// Connection API activation callback to update senders' /transportfile endpoint - captures node_resources by reference!
nmos::connection_sender_transportfile_setter make_node_implementation_transportfile_setter(const nmos::resources& node_resources, const nmos::settings& settings, ConfigManager& config_manager, slog::base_gate& gate)
{
    using web::json::value;

    const auto seed_id = nmos::experimental::fields::seed_id(settings);
    const auto node_id = impl::make_id(seed_id, nmos::types::node);

    auto configIntel = config_manager.get_config();
    auto sender_arr_length = configIntel.senders.size();
    auto sender_arr = configIntel.senders;

    const std::vector<impl::port> media_ports = { impl::ports::video };

    const auto rtp_source_ids = impl::make_ids(seed_id, nmos::types::source, media_ports, sender_arr_length);
    const auto rtp_flow_ids = impl::make_ids(seed_id, nmos::types::flow, media_ports, sender_arr_length);
    const auto rtp_sender_ids = impl::make_ids(seed_id, nmos::types::sender, media_ports, sender_arr_length);

    const uint64_t payload_type_video = impl::fields::sender_payload_type(settings);
    // as part of activation, the example sender /transportfile should be updated based on the active transport parameters
    return [&node_resources, node_id, rtp_source_ids, rtp_flow_ids, rtp_sender_ids, payload_type_video,&gate](const nmos::resource& sender, const nmos::resource& connection_sender, value& endpoint_transportfile)
    {
        const auto found = boost::range::find(rtp_sender_ids, connection_sender.id);
        if (rtp_sender_ids.end() != found)
        {
            const auto index = int(found - rtp_sender_ids.begin());
            const auto source_id = rtp_source_ids.at(index);
            const auto flow_id = rtp_flow_ids.at(index);

            // note, model mutex is already locked by the calling thread, so access to node_resources is OK...
            auto node = nmos::find_resource(node_resources, { node_id, nmos::types::node });
            auto source = nmos::find_resource(node_resources, { source_id, nmos::types::source });
            auto flow = nmos::find_resource(node_resources, { flow_id, nmos::types::flow });
            if (node_resources.end() == node || node_resources.end() == source || node_resources.end() == flow)
            {
                throw std::logic_error("matching IS-04 node, source or flow not found");
            }

            // the nmos::make_sdp_parameters overload from the IS-04 resources provides a high-level interface
            // for common "video/raw", "audio/L", "video/smpte291" and "video/SMPTE2022-6" use cases
            //auto sdp_params = nmos::make_sdp_parameters(node->data, source->data, flow->data, sender.data, { U("PRIMARY"), U("SECONDARY") });

            // nmos::make_{video,audio,data,mux}_sdp_parameters provide a little more flexibility for those four media types
            // and the combination of nmos::make_{video_raw,audio_L,video_smpte291,video_SMPTE2022_6}_parameters
            // with the related make_sdp_parameters overloads provides the most flexible and extensible approach
            auto sdp_params = [&]
            {
                const std::vector<utility::string_t> mids{ U("PRIMARY"), U("SECONDARY") };
                const nmos::format format{ nmos::fields::format(flow->data) };
                if (nmos::formats::video == format)
                {
                    const nmos::media_type video_type{ nmos::fields::media_type(flow->data) };
                    if (nmos::media_types::video_raw == video_type)
                    {
                        return nmos::make_video_sdp_parameters(node->data, source->data, flow->data, sender.data, payload_type_video, mids, {}, sdp::type_parameters::type_N);
                    }
                    else if (nmos::media_types::video_jxsv == video_type)
                    {
                        const auto params = nmos::make_video_jxsv_parameters(node->data, source->data, flow->data, sender.data);
                        const auto ts_refclk = nmos::details::make_ts_refclk(node->data, source->data, sender.data, {});
                        return nmos::make_sdp_parameters(nmos::fields::label(sender.data), params, payload_type_video, mids, ts_refclk);
                    }
                    else
                    {
                        throw std::logic_error("unexpected flow media_type");
                    }
                }
                else if (nmos::formats::audio == format)
                {
                    // this example application doesn't actually stream, so just indicate a sensible value for packet time
                    const double packet_time = nmos::fields::channels(source->data).size() > 8 ? 0.125 : 1;
                    return nmos::make_audio_sdp_parameters(node->data, source->data, flow->data, sender.data, nmos::details::payload_type_audio_default, mids, {}, packet_time);
                }
                else if (nmos::formats::data == format)
                {
                    return nmos::make_data_sdp_parameters(node->data, source->data, flow->data, sender.data, nmos::details::payload_type_data_default, mids, {}, {});
                }
                else if (nmos::formats::mux == format)
                {
                    return nmos::make_mux_sdp_parameters(node->data, source->data, flow->data, sender.data, nmos::details::payload_type_mux_default, mids, {}, sdp::type_parameters::type_N);
                }
                else
                {
                    throw std::logic_error("unexpected flow format");
                }
            }();

            auto& transport_params = nmos::fields::transport_params(nmos::fields::endpoint_active(connection_sender.data));
            auto session_description = nmos::make_session_description(sdp_params, transport_params);
            auto sdp = utility::s2us(sdp::make_session_description(session_description));
            endpoint_transportfile = nmos::make_connection_rtp_sender_transportfile(sdp);
        }
    };
}

// Connection API activation callback to perform application-specific operations to complete activation
nmos::connection_activation_handler make_node_implementation_connection_activation_handler(nmos::node_model& model, ConfigManager& config_manager, AppConnectionResources& app_resources, slog::base_gate& gate)
{
    return [&model, &config_manager, &app_resources, &gate](const nmos::resource& resource, const nmos::resource& connection_resource)
    {
        const std::pair<nmos::id, nmos::type> id_type{ resource.id, resource.type };
        if(id_type.second == nmos::types::sender)
        {
            std::cout<<"Connection API activation handler - sender"<<std::endl;

            // name of environment variable that is set in the sender side linux/container environment
            const char* vfio_port = "VFIO_PORT_TX";
            const char* vfio_port_value_tx = std::getenv(vfio_port);
            auto config_by_id = tracker::get_stream_info(id_type.first);
            if (!vfio_port_value_tx && config_by_id.stream_type.type != stream_type::mcm) {
                // if the stream type is not mcm, then vfio_port_value_tx should be set
                slog::log<slog::severities::error>(gate, SLOG_FLF) << "VFIO_PORT_TX environment variable is not set. You should export one of the virtual function interface port values.";
                return;
            }
            std::thread ffmpegThread1;
            slog::log<slog::severities::info>(gate, SLOG_FLF) << nmos::stash_category(impl::categories::node_implementation) << "this is "<< id_type << "---> sends json for sender";
            auto data = connection_resource.data;
            // Print the full connection_resource data for debugging
            // Extract the sender source IP and receiver destination IP and port from the connection resource data

            auto tp = data[nmos::fields::endpoint_active][nmos::fields::transport_params][0];

            auto sender_source_ip = data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::source_ip];
            auto receiver_destination_ip = data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::destination_ip];
            auto receiver_destination_port = data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::destination_port];

            std::cout << "Sender Source IP: " << sender_source_ip.serialize() << std::endl;
            std::cout << "Receiver Destination IP: " << receiver_destination_ip.serialize() << std::endl;
            std::cout << "Receiver Destination Port: " << receiver_destination_port.serialize() << std::endl;

            // this data are necessary to send via grpc to ffmpeg
            Video v;
            v.frame_width=config_by_id.payload.video.frame_width;
            v.frame_height=config_by_id.payload.video.frame_height;
            v.frame_rate.numerator=config_by_id.payload.video.frame_rate.numerator;
            v.frame_rate.denominator=config_by_id.payload.video.frame_rate.denominator;
            v.pixel_format=config_by_id.payload.video.pixel_format;
            v.video_type=config_by_id.payload.video.video_type;
            v.preset=config_by_id.payload.video.preset;
            v.profile=config_by_id.payload.video.profile;
            Stream s;
            if(config_by_id.stream_type.type==stream_type::mcm)
            {
                s.stream_type.type=stream_type::mcm;
                s.stream_type.mcm.conn_type=config_by_id.stream_type.mcm.conn_type;
                s.stream_type.mcm.transport=config_by_id.stream_type.mcm.transport;
                s.stream_type.mcm.transport_pixel_format=config_by_id.stream_type.mcm.transport_pixel_format;
                s.stream_type.mcm.ip=receiver_destination_ip.as_string();
                s.stream_type.mcm.port=receiver_destination_port.as_integer();
                s.stream_type.mcm.urn=config_by_id.stream_type.mcm.urn;
            }
            else if(config_by_id.stream_type.type==stream_type::st2110)
            {
                s.stream_type.type = stream_type::st2110;
                s.stream_type.st2110.network_interface = vfio_port_value_tx;
                s.stream_type.st2110.local_ip = sender_source_ip.as_string();
                s.stream_type.st2110.remote_ip = receiver_destination_ip.as_string();
                s.stream_type.st2110.transport = config_by_id.stream_type.st2110.transport;
                s.stream_type.st2110.remote_port = receiver_destination_port.as_integer();
                s.stream_type.st2110.payload_type = impl::fields::sender_payload_type(model.settings);
                s.stream_type.st2110.queues_cnt = config_by_id.stream_type.st2110.queues_cnt;
            }

            Payload payload;
            payload.type = payload_type::video;
            payload.video = v;

            s.payload = payload;

            // sender=nmos, receiver=ffmpeg
            // to get ffmpeg receivers of stream_type::File
            auto configIntel = config_manager.get_config();
            config_manager.print_config();

            auto ffmpeg_receiver_as_file_vector = tracker::get_file_streams_receivers(configIntel);
            for (auto& stream_receiver : ffmpeg_receiver_as_file_vector) {
                stream_receiver.payload.type = payload_type::video;
            }

            auto gpu_hw_acceleration_device = "";
            if (configIntel.gpu_hw_acceleration == "intel") {
                if (configIntel.gpu_hw_acceleration_device.empty()) {
                    slog::log<slog::severities::error>(gate, SLOG_FLF) << "GPU hardware acceleration device is not specified for Intel.";
                }
                gpu_hw_acceleration_device = configIntel.gpu_hw_acceleration_device.c_str();
            }
            // construct config for NMOS sender

            if ( app_resources.all_senders.size() < configIntel.senders.size()) {
                app_resources.all_senders.push_back(s);
                slog::log<slog::severities::info>(gate, SLOG_FLF) << "New sender added to the list. Total senders: " << app_resources.all_senders.size();
            }

            const Config config = {app_resources.all_senders, ffmpeg_receiver_as_file_vector, configIntel.function, configIntel.multiviewer_columns, configIntel.gpu_hw_acceleration, gpu_hw_acceleration_device, configIntel.stream_loop, configIntel.logging_level};

            if ( app_resources.all_senders.size() < configIntel.senders.size()) {
                slog::log<slog::severities::info>(gate, SLOG_FLF) <<  "Waiting for connection to all declared senders. Connected " << \
                app_resources.all_senders.size() << " senders from declared " << configIntel.senders.size() << " senders";
            }else{
                ffmpegThread1=std::thread(grpc::sendDataToFfmpeg, impl::fields::ffmpeg_grpc_server_address(model.settings), impl::fields::ffmpeg_grpc_server_port(model.settings), config);
                app_resources.threads.push_back(std::move(ffmpegThread1));
                app_resources.all_senders.clear();
            }
}
        if(id_type.second == nmos::types::receiver)
        {

            // Right now only connection IP addresses (interface_ip, multicast_ip and source_ip)
            // and destination_port are only overidden by transport_params.
            // This means that transportfile SDP has lower priority than transport_params data.

            std::cout<<"Connection API activation handler - receiver"<<std::endl;
            std::thread ffmpegThread2;

            // name of environment variable that is set in the receiver side linux/container environment
            const char* vfio_port = "VFIO_PORT_RX";
            const char* vfio_port_value_rx = std::getenv(vfio_port);

            // get the stream type
            auto config_by_id = tracker::get_stream_info(id_type.first);
            if (!vfio_port_value_rx && config_by_id.stream_type.type != stream_type::mcm) {
                // if the stream type is not mcm, then vfio_port_value_rx should be set
                slog::log<slog::severities::error>(gate, SLOG_FLF) << "VFIO_PORT_RX environment variable is not set. You should export one of the virtual function interface port values.";
                return;
            }

            //retrieve the receiver configuration data
            auto data = connection_resource.data;

            //interface_ip MUST be set on receiver side
            if (data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::interface_ip].is_null()) {
                slog::log<slog::severities::error>(gate, SLOG_FLF) << "interface_ip is null. Please set the interface_ip in the transport_params on receiver NMOS side.";
                return;
            }

            // Print the transport parameters that are parsed based on REST API request
            std::cout << "Receiver ID: " << id_type.first << std::endl;
            std::cout << "TRANSPORT PARAMS: " <<  std::endl;
            std::cout << "Receiver Interface IP: " << data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::interface_ip].serialize() << std::endl;
            std::cout << "Receiver Multicast IP: " << data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::multicast_ip].serialize() << std::endl;
            std::cout << "Receiver Source IP: " << data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::source_ip].serialize() << std::endl;
            std::cout << "Receiver Destination Port: " << data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::destination_port].serialize() << std::endl;
            std::cout << "RTP Enabled: " << data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::rtp_enabled].serialize() << std::endl;

            auto trasportfile_sdp = data[nmos::fields::endpoint_active][nmos::fields::transport_file][nmos::fields::data];
            std::cout<<utility::us2s(trasportfile_sdp.as_string());
            auto session_description = sdp::parse_session_description(utility::us2s(trasportfile_sdp.as_string()));
            auto transport_params_from_sdp = nmos::get_session_description_transport_params(session_description);

            // Print the transport parameters that are parsed based on SDP file
            std::cout << "SDP Transport Parameter: destination_port:  " << transport_params_from_sdp[0][nmos::fields::destination_port]<< std::endl;
            std::cout << "SDP Transport Parameter: interface_ip:  " << transport_params_from_sdp[0][nmos::fields::interface_ip]<< std::endl;
            std::cout << "SDP Transport Parameter: multicast_ip:  " << transport_params_from_sdp[0][nmos::fields::multicast_ip]<< std::endl;
            std::cout << "SDP Transport Parameter: source_ip:  " << transport_params_from_sdp[0][nmos::fields::source_ip]<< std::endl;
            std::cout << "SDP Transport Parameter: rtp_enabled:  " << transport_params_from_sdp[0][nmos::fields::rtp_enabled]<< std::endl;

            auto rx_destination_port = transport_params_from_sdp[0][nmos::fields::destination_port];
            auto rx_interface_ip = transport_params_from_sdp[0][nmos::fields::interface_ip];
            auto rx_multicast_ip = transport_params_from_sdp[0][nmos::fields::multicast_ip];
            auto rx_source_ip = transport_params_from_sdp[0][nmos::fields::source_ip];
            
            if (rx_interface_ip == transport_params_from_sdp[0][nmos::fields::interface_ip] &&
                rx_multicast_ip == transport_params_from_sdp[0][nmos::fields::multicast_ip] &&
                rx_source_ip == transport_params_from_sdp[0][nmos::fields::source_ip] &&
                rx_destination_port == transport_params_from_sdp[0][nmos::fields::destination_port])
                {
                std::cout << "The transport parameters from the active connection and the SDP are equal. Using parameters from SDP transportfile" << std::endl;
            } else {
                std::cout << "The transport parameters from the active connection and the SDP are not equal. Using parameters from REST API request transport_params" << std::endl;
                rx_destination_port = data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::destination_port];
                rx_interface_ip = data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::interface_ip];
                rx_multicast_ip = data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::multicast_ip];
                rx_source_ip = data[nmos::fields::endpoint_active][nmos::fields::transport_params][0][nmos::fields::source_ip];
            }

            auto& media_descriptions = sdp::fields::media_descriptions(session_description);
            auto& media_description = media_descriptions.at(0);
            auto& media = sdp::fields::media(media_description);
            auto& attributes = sdp::fields::attributes(media_description).as_array();
            auto rtpmap = sdp::find_name(attributes, sdp::attributes::rtpmap);
            auto& encoding_name = sdp::fields::encoding_name(sdp::fields::value(*rtpmap));
            const auto& payload_type_sdp = sdp::fields::payload_type(sdp::fields::value(*rtpmap));
            auto source_filter = sdp::find_name(attributes, sdp::attributes::source_filter);
            auto fmtp = sdp::find_name(attributes, sdp::attributes::fmtp);
            auto& params = sdp::fields::format_specific_parameters(sdp::fields::value(*fmtp));
            auto width_param = sdp::find_name(params, U("width"));
            auto height_param = sdp::find_name(params, U("height"));
            auto fps = sdp::find_name(params, U("exactframerate"));

            auto fps_numerator = nmos::details::parse_exactframerate(sdp::fields::value(*fps).as_string()).numerator();
            auto fps_denominator = nmos::details::parse_exactframerate(sdp::fields::value(*fps).as_string()).denominator();

            // this data are necessary to send via grpc to ffmpeg
            config_manager.print_config();

            Video v;
            v.frame_width=std::stoi(sdp::fields::value(*width_param).as_string());
            v.frame_height=std::stoi(sdp::fields::value(*height_param).as_string());
            v.frame_rate.numerator=fps_numerator;
            v.frame_rate.denominator=fps_denominator;
            v.pixel_format=config_by_id.payload.video.pixel_format;
            v.preset=config_by_id.payload.video.preset;
            v.profile=config_by_id.payload.video.profile;
            if (encoding_name == U("raw")) {
                v.video_type = "rawvideo";
            } else {
                v.video_type = utility::us2s(encoding_name);
            }
            Stream s;
            if(config_by_id.stream_type.type==stream_type::mcm)
            {
                s.stream_type.type=stream_type::mcm;
                s.stream_type.mcm.conn_type=config_by_id.stream_type.mcm.conn_type;
                s.stream_type.mcm.transport=config_by_id.stream_type.mcm.transport;
                s.stream_type.mcm.transport_pixel_format=config_by_id.stream_type.mcm.transport_pixel_format;
                if (rx_multicast_ip.is_null()) {
                    s.stream_type.mcm.ip = rx_source_ip.as_string();
                } else {
                    s.stream_type.mcm.ip = rx_multicast_ip.as_string();
                }
                s.stream_type.mcm.port=rx_destination_port.as_integer();
                s.stream_type.mcm.urn=config_by_id.stream_type.mcm.urn;
            }
            else if(config_by_id.stream_type.type==stream_type::st2110)
            {
                s.stream_type.type = stream_type::st2110;
                s.stream_type.st2110.network_interface = vfio_port_value_rx;
                s.stream_type.st2110.local_ip = rx_interface_ip.as_string();
                if (rx_multicast_ip.is_null()) {
                    s.stream_type.st2110.remote_ip = rx_source_ip.as_string();
                } else {
                    s.stream_type.st2110.remote_ip = rx_multicast_ip.as_string();
                }
                
                s.stream_type.st2110.transport = config_by_id.stream_type.st2110.transport;
                s.stream_type.st2110.remote_port = rx_destination_port.as_integer();
                s.stream_type.st2110.payload_type = payload_type_sdp;
                s.stream_type.st2110.queues_cnt = config_by_id.stream_type.st2110.queues_cnt;
            }

            Payload payload;
            payload.type = payload_type::video;
            payload.video = v;

            s.payload = payload;

            //receiver=nmos, sender=ffmpeg
            //to get ffmpeg senders of stream_type::File
            auto configIntel = config_manager.get_config();
            auto ffmpeg_sender_as_file_vector = tracker::get_file_streams_senders(configIntel);
            for (auto& stream_sender : ffmpeg_sender_as_file_vector) {
                stream_sender.payload.type = payload_type::video;
            }
            if (ffmpeg_sender_as_file_vector.empty()) {
                slog::log<slog::severities::error>(gate, SLOG_FLF) << "No ffmpeg senders of stream_type::File found.";
            }
            auto gpu_hw_acceleration_device = "";
            if (configIntel.gpu_hw_acceleration == "intel") {
                if (configIntel.gpu_hw_acceleration_device.empty()) {
                    slog::log<slog::severities::error>(gate, SLOG_FLF) << "GPU hardware acceleration device is not specified for Intel.";
                }
                gpu_hw_acceleration_device = configIntel.gpu_hw_acceleration_device.c_str();
            }
            // construct config for NMOS sender
            if ( app_resources.all_receivers.size() < configIntel.receivers.size()) {
                app_resources.all_receivers.push_back(s);
                slog::log<slog::severities::info>(gate, SLOG_FLF) << "New receiver added to the list. Total receivers: " << app_resources.all_receivers.size();
            }
            const Config config = {ffmpeg_sender_as_file_vector,app_resources.all_receivers,configIntel.function,configIntel.multiviewer_columns, configIntel.gpu_hw_acceleration,gpu_hw_acceleration_device, configIntel.stream_loop, configIntel.logging_level};

            if ( app_resources.all_receivers.size() < configIntel.receivers.size()) {
                slog::log<slog::severities::info>(gate, SLOG_FLF) <<  "Waiting for connection to all declared receivers. Connected " << \
                app_resources.all_receivers.size() << " receivers from declared " << configIntel.receivers.size() << " receivers";
            }else{
                ffmpegThread2=std::thread(grpc::sendDataToFfmpeg, impl::fields::ffmpeg_grpc_server_address(model.settings), \
                impl::fields::ffmpeg_grpc_server_port(model.settings), config);
                app_resources.threads.push_back(std::move(ffmpegThread2));
                app_resources.all_receivers.clear();
            }
        }
    };
}

namespace impl
{
    nmos::interlace_mode get_interlace_mode(const nmos::rational& frame_rate, uint32_t frame_height, const nmos::settings& settings)
    {
        if (settings.has_field(impl::fields::interlace_mode))
        {
            return nmos::interlace_mode{ impl::fields::interlace_mode(settings) };
        }
        // for the default, 1080i50 and 1080i59.94 are arbitrarily preferred to 1080p25 and 1080p29.97
        // for 1080i formats, ST 2110-20 says that "the fields of an interlaced image are transmitted in time order,
        // first field first [and] the sample rows of the temporally second field are displaced vertically 'below' the
        // like-numbered sample rows of the temporally first field."
        // const auto frame_rate = nmos::parse_rational(impl::fields::frame_rate(settings));
        // const auto frame_height = impl::fields::frame_height(settings);
        return (nmos::rates::rate25 == frame_rate || nmos::rates::rate29_97 == frame_rate) && 1080 == frame_height
            ? nmos::interlace_modes::interlaced_tff
            : nmos::interlace_modes::progressive;
    }

    // find interface with the specified address
    std::vector<web::hosts::experimental::host_interface>::const_iterator find_interface(const std::vector<web::hosts::experimental::host_interface>& interfaces, const utility::string_t& address)
    {
        return boost::range::find_if(interfaces, [&](const web::hosts::experimental::host_interface& interface)
        {
            return interface.addresses.end() != boost::range::find(interface.addresses, address);
        });
    }

    // generate repeatable ids for the example node's resources
    nmos::id make_id(const nmos::id& seed_id, const nmos::type& type, const impl::port& port, int index)
    {
        return nmos::make_repeatable_id(seed_id, U("/x-nmos/node/") + type.name + U('/') + port.name + utility::conversions::details::to_string_t(index));
    }

    std::vector<nmos::id> make_ids(const nmos::id& seed_id, const nmos::type& type, const impl::port& port, int how_many)
    {
        return boost::copy_range<std::vector<nmos::id>>(boost::irange(0, how_many) | boost::adaptors::transformed([&](const int& index)
        {
            return impl::make_id(seed_id, type, port, index);
        }));
    }

    std::vector<nmos::id> make_ids(const nmos::id& seed_id, const nmos::type& type, const std::vector<port>& ports, int how_many)
    {
        // hm, boost::range::combine arrived in Boost 1.56.0
        std::vector<nmos::id> ids;
        for (const auto& port : ports)
        {
            boost::range::push_back(ids, make_ids(seed_id, type, port, how_many));
        }
        return ids;
    }

    std::vector<nmos::id> make_ids(const nmos::id& seed_id, const std::vector<nmos::type>& types, const std::vector<port>& ports, int how_many)
    {
        // hm, boost::range::combine arrived in Boost 1.56.0
        std::vector<nmos::id> ids;
        for (const auto& type : types)
        {
            boost::range::push_back(ids, make_ids(seed_id, type, ports, how_many));
        }
        return ids;
    }

    // generate a repeatable source-specific multicast address for each leg of a sender
    utility::string_t make_source_specific_multicast_address_v4(const nmos::id& id, int leg)
    {
        // hash the pseudo-random id and leg to generate the address
        const auto s = id + U('/') + utility::conversions::details::to_string_t(leg);
        const auto h = std::hash<utility::string_t>{}(s);
        auto a = boost::asio::ip::address_v4(uint32_t(h)).to_bytes();
        // ensure the address is in the source-specific multicast block reserved for local host allocation, 232.0.1.0-232.255.255.255
        // see https://www.iana.org/assignments/multicast-addresses/multicast-addresses.xhtml#multicast-addresses-10
        a[0] = 232;
        a[2] |= 1;
        return utility::s2us(boost::asio::ip::address_v4(a).to_string());
    }

    // add a selection of parents to a source or flow
    void insert_parents(nmos::resource& resource, const nmos::id& seed_id, const port& port, int index)
    {
        // algorithm to produce signal ancestry with a range of depths and breadths
        // see https://github.com/sony/nmos-cpp/issues/312#issuecomment-1335641637
        int b = 0;
        while (index & (1 << b)) ++b;
        if (!b) return;
        index &= ~(1 << (b - 1));
        do
        {
            index &= ~(1 << b);
            web::json::push_back(resource.data[nmos::fields::parents], impl::make_id(seed_id, resource.type, port, index));
            ++b;
        } while (index & (1 << b));
    }

    // add a helpful suffix to the label of a sub-resource for the example node
    void set_label_description(nmos::resource& resource, const impl::port& port, int index)
    {
        using web::json::value;

        auto label = nmos::fields::label(resource.data);
        if (!label.empty()) label += U('/');
        label += resource.type.name + U('/') + port.name + utility::conversions::details::to_string_t(index);
        resource.data[nmos::fields::label] = value::string(label);

        auto description = nmos::fields::description(resource.data);
        if (!description.empty()) description += U('/');
        description += resource.type.name + U('/') + port.name + utility::conversions::details::to_string_t(index);
        resource.data[nmos::fields::description] = value::string(description);
    }

    // add an example "natural grouping" hint to a sender or receiver
    void insert_group_hint(nmos::resource& resource, const impl::port& port, int index)
    {
        web::json::push_back(resource.data[nmos::fields::tags][nmos::fields::group_hint], nmos::make_group_hint({ U("example"), resource.type.name + U(' ') + port.name + utility::conversions::details::to_string_t(index) }));
    }
}

// This constructs all the callbacks used to integrate the example device-specific underlying implementation
// into the server instance for the NMOS Node.
nmos::experimental::node_implementation make_node_implementation(nmos::node_model& model, ConfigManager& config_manager,AppConnectionResources& app_resources, slog::base_gate& gate)
{
    return nmos::experimental::node_implementation()
        .on_load_server_certificates(nmos::make_load_server_certificates_handler(model.settings, gate))
        .on_load_dh_param(nmos::make_load_dh_param_handler(model.settings, gate))
        .on_load_ca_certificates(nmos::make_load_ca_certificates_handler(model.settings, gate))
        .on_system_changed(make_node_implementation_system_global_handler(model, gate)) // may be omitted if not required
        .on_registration_changed(make_node_implementation_registration_handler(gate)) // may be omitted if not required
        .on_parse_transport_file(make_node_implementation_transport_file_parser(gate)) // may be omitted if the default is sufficient
        .on_validate_connection_resource_patch(make_node_implementation_patch_validator(gate)) // may be omitted if not required
        .on_resolve_auto(make_node_implementation_auto_resolver(model.settings, config_manager, gate))
        .on_set_transportfile(make_node_implementation_transportfile_setter(model.node_resources, model.settings, config_manager, gate))
        .on_connection_activated(make_node_implementation_connection_activation_handler(model, config_manager, app_resources, gate));
}
