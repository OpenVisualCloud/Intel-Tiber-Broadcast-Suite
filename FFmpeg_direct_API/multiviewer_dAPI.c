#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavfilter/avfilter.h>
#include <libavfilter/buffersink.h>
#include <libavfilter/buffersrc.h>
#include <libavutil/opt.h>
#include <libavutil/imgutils.h>
#include <libavutil/hwcontext_qsv.h>
#include <libswscale/swscale.h>

#define NUM_INPUTS 8

// Function to set up the filter graph
int setup_filter_graph(AVFilterGraph **graph, AVFilterContext **buffersrc_ctx, AVFilterContext **buffersink_ctx, AVCodecContext *dec_ctx[], int num_inputs) {
    AVFilterGraph *filter_graph = avfilter_graph_alloc();
    if (!filter_graph) {
        fprintf(stderr, "Could not allocate filter graph\n");
        return -1;
    }

    char args[512];
    AVFilterContext *hwupload_ctx[num_inputs];
    AVFilterContext *scale_ctx[num_inputs];
    AVFilterContext *xstack_ctx;
    const AVFilter *buffersrc = avfilter_get_by_name("buffer");
    const AVFilter *buffersink = avfilter_get_by_name("buffersink");
    const AVFilter *hwupload = avfilter_get_by_name("hwupload");
    const AVFilter *scale_qsv = avfilter_get_by_name("scale_qsv");
    const AVFilter *xstack_qsv = avfilter_get_by_name("xstack_qsv");
    const AVFilter *format = avfilter_get_by_name("format");

    // Create buffer source filters
    for (int i = 0; i < num_inputs; i++) {
        snprintf(args, sizeof(args),
                 "video_size=%dx%d:pix_fmt=%d:time_base=%d/%d:pixel_aspect=%d/%d",
                 dec_ctx[i]->width, dec_ctx[i]->height, dec_ctx[i]->pix_fmt,
                 dec_ctx[i]->time_base.num, dec_ctx[i]->time_base.den,
                 dec_ctx[i]->sample_aspect_ratio.num, dec_ctx[i]->sample_aspect_ratio.den);
        if (avfilter_graph_create_filter(&buffersrc_ctx[i], buffersrc, NULL, args, NULL, filter_graph) < 0) {
            fprintf(stderr, "Could not create buffer source filter for input %d\n", i);
            avfilter_graph_free(&filter_graph);
            return -1;
        }
    }

    // Create hwupload and scale filters
    for (int i = 0; i < num_inputs; i++) {
        if (avfilter_graph_create_filter(&hwupload_ctx[i], hwupload, NULL, NULL, NULL, filter_graph) < 0) {
            fprintf(stderr, "Could not create hwupload filter for input %d\n", i);
            avfilter_graph_free(&filter_graph);
            return -1;
        }
        snprintf(args, sizeof(args), "iw/4:ih/2");
        if (avfilter_graph_create_filter(&scale_ctx[i], scale_qsv, NULL, args, NULL, filter_graph) < 0) {
            fprintf(stderr, "Could not create scale_qsv filter for input %d\n", i);
            avfilter_graph_free(&filter_graph);
            return -1;
        }
    }

    // Create xstack filter
    if (avfilter_graph_create_filter(&xstack_ctx, xstack_qsv, NULL, "inputs=8:layout=0_0|w0_0|0_h0|w0_h0|w0+w1_0|w0+w1+w2_0|w0+w1_h0|w0+w1+w2_h0", NULL, filter_graph) < 0) {
        fprintf(stderr, "Could not create xstack_qsv filter\n");
        avfilter_graph_free(&filter_graph);
        return -1;
    }

    // Create format filters
    AVFilterContext *format_ctx1;
    if (avfilter_graph_create_filter(&format_ctx1, format, NULL, "y210le", NULL, filter_graph) < 0) {
        fprintf(stderr, "Could not create format filter y210le\n");
        avfilter_graph_free(&filter_graph);
        return -1;
    }
    AVFilterContext *format_ctx2;
    if (avfilter_graph_create_filter(&format_ctx2, format, NULL, "yuv422p10le", NULL, filter_graph) < 0) {
        fprintf(stderr, "Could not create format filter yuv422p10le\n");
        avfilter_graph_free(&filter_graph);
        return -1;
    }

    // Create buffer sink filter
    if (avfilter_graph_create_filter(buffersink_ctx, buffersink, NULL, NULL, NULL, filter_graph) < 0) {
        fprintf(stderr, "Could not create buffer sink filter\n");
        avfilter_graph_free(&filter_graph);
        return -1;
    }

    // Link filters
    for (int i = 0; i < num_inputs; i++) {
        if (avfilter_link(buffersrc_ctx[i], 0, hwupload_ctx[i], 0) < 0 ||
            avfilter_link(hwupload_ctx[i], 0, scale_ctx[i], 0) < 0 ||
            avfilter_link(scale_ctx[i], 0, xstack_ctx, i) < 0) {
            fprintf(stderr, "Error linking filters for input %d\n", i);
            avfilter_graph_free(&filter_graph);
            return -1;
        }
    }
    if (avfilter_link(xstack_ctx, 0, format_ctx1, 0) < 0 ||
        avfilter_link(format_ctx1, 0, format_ctx2, 0) < 0 ||
        avfilter_link(format_ctx2, 0, *buffersink_ctx, 0) < 0) {
        fprintf(stderr, "Error linking xstack and format filters\n");
        avfilter_graph_free(&filter_graph);
        return -1;
    }

    // Configure the filter graph
    if (avfilter_graph_config(filter_graph, NULL) < 0) {
        fprintf(stderr, "Error configuring the filter graph\n");
        avfilter_graph_free(&filter_graph);
        return -1;
    }

    *graph = filter_graph;
    return 0;
}

// Function to start processing
void start_processing(const char *input_filenames[], int num_inputs, const char *output_filename) {
    AVFormatContext *input_fmt_ctx[num_inputs];
    AVFormatContext *output_fmt_ctx = NULL;
    AVCodecContext *dec_ctx[num_inputs];
    AVCodecContext *enc_ctx = NULL;
    AVFilterGraph *filter_graph = NULL;
    AVFilterContext *buffersrc_ctx[num_inputs];
    AVFilterContext *buffersink_ctx = NULL;
    int ret;

    // Open input files
    for (int i = 0; i < num_inputs; i++) {
        AVDictionary *options = NULL;
        av_dict_set(&options, "p_port", "${VFIO_PORT_PROC}", 0);
        av_dict_set(&options, "p_sip", "192.168.2.2", 0);
        av_dict_set(&options, "p_rx_ip", "192.168.2.1", 0);
        av_dict_set(&options, "udp_port", "20000", 0);
        av_dict_set(&options, "payload_type", "112", 0);
        av_dict_set(&options, "fps", "25", 0);
        av_dict_set(&options, "pix_fmt", "yuv422p10le", 0);
        av_dict_set(&options, "video_size", "1920x1080", 0);

        if (avformat_open_input(&input_fmt_ctx[i], input_filenames[i], av_find_input_format("mtl_st20p"), &options) < 0) {
            fprintf(stderr, "Could not open input file %s\n", input_filenames[i]);
            av_dict_free(&options);
            goto end;
        }
        av_dict_free(&options);

        // Find stream info
        if (avformat_find_stream_info(input_fmt_ctx[i], NULL) < 0) {
            fprintf(stderr, "Could not find stream info for input file %s\n", input_filenames[i]);
            goto end;
        }

        // Find video stream
        int video_stream_index = av_find_best_stream(input_fmt_ctx[i], AVMEDIA_TYPE_VIDEO, -1, -1, NULL, 0);
        if (video_stream_index < 0) {
            fprintf(stderr, "Could not find video stream in input file %s\n", input_filenames[i]);
            goto end;
        }

        // Set up decoder
        AVStream *video_stream = input_fmt_ctx[i]->streams[video_stream_index];
        AVCodec *dec = avcodec_find_decoder(video_stream->codecpar->codec_id);
        if (!dec) {
            fprintf(stderr, "Could not find decoder for input file %s\n", input_filenames[i]);
            goto end;
        }
        dec_ctx[i] = avcodec_alloc_context3(dec);
        if (!dec_ctx[i]) {
            fprintf(stderr, "Could not allocate decoder context for input file %s\n", input_filenames[i]);
            goto end;
        }
        if (avcodec_parameters_to_context(dec_ctx[i], video_stream->codecpar) < 0) {
            fprintf(stderr, "Could not copy codec parameters to context for input file %s\n", input_filenames[i]);
            goto end;
        }
        if (avcodec_open2(dec_ctx[i], dec, NULL) < 0) {
            fprintf(stderr, "Could not open decoder for input file %s\n", input_filenames[i]);
            goto end;
        }
    }

    // Set up filter graph
    if (setup_filter_graph(&filter_graph, buffersrc_ctx, &buffersink_ctx, dec_ctx, num_inputs) < 0) {
        fprintf(stderr, "Could not set up filter graph\n");
        goto end;
    }

    // Set up output format context
    AVDictionary *output_options = NULL;
    av_dict_set(&output_options, "p_port", "${VFIO_PORT_PROC}", 0);
    av_dict_set(&output_options, "p_sip", "192.168.2.2", 0);
    av_dict_set(&output_options, "p_tx_ip", "192.168.2.3", 0);
    av_dict_set(&output_options, "udp_port", "20000", 0);
    av_dict_set(&output_options, "payload_type", "112", 0);
    av_dict_set(&output_options, "fps", "25", 0);
    av_dict_set(&output_options, "pix_fmt", "yuv422p10le", 0);

    avformat_alloc_output_context2(&output_fmt_ctx, av_find_output_format("mtl_st20p"), NULL, output_filename);
    if (!output_fmt_ctx) {
        fprintf(stderr, "Could not create output context\n");
        av_dict_free(&output_options);
        goto end;
    }
    AVOutputFormat *output_fmt = output_fmt_ctx->oformat;

    // Create output stream
    AVStream *out_stream = avformat_new_stream(output_fmt_ctx, NULL);
    if (!out_stream) {
        fprintf(stderr, "Could not create output stream\n");
        av_dict_free(&output_options);
        goto end;
    }

    // Set up encoder
    AVCodec *enc = avcodec_find_encoder(AV_CODEC_ID_H264);
    if (!enc) {
        fprintf(stderr, "Could not find encoder\n");
        av_dict_free(&output_options);
        goto end;
    }
    enc_ctx = avcodec_alloc_context3(enc);
    if (!enc_ctx) {
        fprintf(stderr, "Could not allocate encoder context\n");
        av_dict_free(&output_options);
        goto end;
    }
    enc_ctx->height = dec_ctx[0]->height / 2; // Because of scale_qsv
    enc_ctx->width = dec_ctx[0]->width / 4 * 4; // Because of scale_qsv and xstack
    enc_ctx->sample_aspect_ratio = dec_ctx[0]->sample_aspect_ratio;
    enc_ctx->pix_fmt = AV_PIX_FMT_YUV422P10LE;
    enc_ctx->time_base = (AVRational){1, 25};
    if (output_fmt->flags & AVFMT_GLOBALHEADER) {
        enc_ctx->flags |= AV_CODEC_FLAG_GLOBAL_HEADER;
    }
    if (avcodec_open2(enc_ctx, enc, NULL) < 0) {
        fprintf(stderr, "Could not open encoder\n");
        av_dict_free(&output_options);
        goto end;
    }
    if (avcodec_parameters_from_context(out_stream->codecpar, enc_ctx) < 0) {
        fprintf(stderr, "Could not copy encoder parameters to stream\n");
        av_dict_free(&output_options);
        goto end;
    }

    // Open output file
    if (!(output_fmt->flags & AVFMT_NOFILE)) {
        if (avio_open2(&output_fmt_ctx->pb, output_filename, AVIO_FLAG_WRITE, NULL, &output_options) < 0) {
            fprintf(stderr, "Could not open output file\n");
            av_dict_free(&output_options);
            goto end;
        }
    }
    av_dict_free(&output_options);

    // Write output file header
    if (avformat_write_header(output_fmt_ctx, NULL) < 0) {
        fprintf(stderr, "Could not write output file header\n");
        goto end;
    }

    // Allocate frames and packets
    AVFrame *frame = av_frame_alloc();
    AVFrame *filt_frame = av_frame_alloc();
    AVPacket *pkt = av_packet_alloc();

    // Read, decode, filter, encode, and write frames
    for (int i = 0; i < num_inputs; i++) {
        while (av_read_frame(input_fmt_ctx[i], pkt) >= 0) {
            if (pkt->stream_index == 0) {
                ret = avcodec_send_packet(dec_ctx[i], pkt);
                if (ret < 0) {
                    fprintf(stderr, "Error sending packet to decoder for input file %s\n", input_filenames[i]);
                    break;
                }
                while (ret >= 0) {
                    ret = avcodec_receive_frame(dec_ctx[i], frame);
                    if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) {
                        break;
                    } else if (ret < 0) {
                        fprintf(stderr, "Error receiving frame from decoder for input file %s\n", input_filenames[i]);
                        goto end;
                    }

                    // Send frame to filter graph
                    if (av_buffersrc_add_frame(buffersrc_ctx[i], frame) < 0) {
                        fprintf(stderr, "Error sending frame to filter graph for input file %s\n", input_filenames[i]);
                        goto end;
                    }

                    // Receive filtered frame
                    while (1) {
                        ret = av_buffersink_get_frame(buffersink_ctx, filt_frame);
                        if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) {
                            break;
                        } else if (ret < 0) {
                            fprintf(stderr, "Error receiving frame from filter graph\n");
                            goto end;
                        }

                        // Encode filtered frame
                        ret = avcodec_send_frame(enc_ctx, filt_frame);
                        if (ret < 0) {
                            fprintf(stderr, "Error sending frame to encoder\n");
                            goto end;
                        }
                        while (ret >= 0) {
                            ret = avcodec_receive_packet(enc_ctx, pkt);
                            if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) {
                                break;
                            } else if (ret < 0) {
                                fprintf(stderr, "Error receiving packet from encoder\n");
                                goto end;
                            }

                            // Write packet to output file
                            pkt->stream_index = out_stream->index;
                            av_packet_rescale_ts(pkt, enc_ctx->time_base, out_stream->time_base);
                            if (av_interleaved_write_frame(output_fmt_ctx, pkt) < 0) {
                                fprintf(stderr, "Error writing packet to output file\n");
                                goto end;
                            }
                            av_packet_unref(pkt);
                        }
                        av_frame_unref(filt_frame);
                    }
                    av_frame_unref(frame);
                }
            }
            av_packet_unref(pkt);
        }
    }

    // Write output file trailer
    av_write_trailer(output_fmt_ctx);

end:
    // Cleanup
    av_frame_free(&frame);
    av_frame_free(&filt_frame);
    av_packet_free(&pkt);
    for (int i = 0; i < num_inputs; i++) {
        avcodec_free_context(&dec_ctx[i]);
        avformat_close_input(&input_fmt_ctx[i]);
    }
    avcodec_free_context(&enc_ctx);
    if (output_fmt_ctx && !(output_fmt_ctx->oformat->flags & AVFMT_NOFILE)) {
        avio_closep(&output_fmt_ctx->pb);
    }
    avformat_free_context(output_fmt_ctx);
    if (filter_graph) {
        avfilter_graph_free(&filter_graph);
    }
}

#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavfilter/avfilter.h>
#include <libavfilter/buffersink.h>
#include <libavfilter/buffersrc.h>
#include <libavutil/opt.h>
#include <libavutil/imgutils.h>
#include <libavutil/hwcontext_qsv.h>
#include <libswscale/swscale.h>

#define NUM_INPUTS 8

// Function to initialize FFmpeg libraries on load
__attribute__((constructor))
static void init_ffmpeg() {
    
}

// Function to set up the filter graph
int setup_filter_graph(AVFilterGraph **graph, AVFilterContext **buffersrc_ctx, AVFilterContext **buffersink_ctx, AVCodecContext *dec_ctx[], int num_inputs) {
    AVFilterGraph *filter_graph = avfilter_graph_alloc();
    if (!filter_graph) {
        fprintf(stderr, "Could not allocate filter graph\n");
        return -1;
    }

    char args[512];
    AVFilterContext *hwupload_ctx[num_inputs];
    AVFilterContext *scale_ctx[num_inputs];
    AVFilterContext *xstack_ctx;
    const AVFilter *buffersrc = avfilter_get_by_name("buffer");
    const AVFilter *buffersink = avfilter_get_by_name("buffersink");
    const AVFilter *hwupload = avfilter_get_by_name("hwupload");
    const AVFilter *scale_qsv = avfilter_get_by_name("scale_qsv");
    const AVFilter *xstack_qsv = avfilter_get_by_name("xstack_qsv");
    const AVFilter *format = avfilter_get_by_name("format");

    // Create buffer source filters
    for (int i = 0; i < num_inputs; i++) {
        snprintf(args, sizeof(args),
                 "video_size=%dx%d:pix_fmt=%d:time_base=%d/%d:pixel_aspect=%d/%d",
                 dec_ctx[i]->width, dec_ctx[i]->height, dec_ctx[i]->pix_fmt,
                 dec_ctx[i]->time_base.num, dec_ctx[i]->time_base.den,
                 dec_ctx[i]->sample_aspect_ratio.num, dec_ctx[i]->sample_aspect_ratio.den);
        if (avfilter_graph_create_filter(&buffersrc_ctx[i], buffersrc, NULL, args, NULL, filter_graph) < 0) {
            fprintf(stderr, "Could not create buffer source filter for input %d\n", i);
            avfilter_graph_free(&filter_graph);
            return -1;
        }
    }

    // Create hwupload and scale filters
    for (int i = 0; i < num_inputs; i++) {
        if (avfilter_graph_create_filter(&hwupload_ctx[i], hwupload, NULL, NULL, NULL, filter_graph) < 0) {
            fprintf(stderr, "Could not create hwupload filter for input %d\n", i);
            avfilter_graph_free(&filter_graph);
            return -1;
        }
        snprintf(args, sizeof(args), "iw/4:ih/2");
        if (avfilter_graph_create_filter(&scale_ctx[i], scale_qsv, NULL, args, NULL, filter_graph) < 0) {
            fprintf(stderr, "Could not create scale_qsv filter for input %d\n", i);
            avfilter_graph_free(&filter_graph);
            return -1;
        }
    }

    // Create xstack filter
    if (avfilter_graph_create_filter(&xstack_ctx, xstack_qsv, NULL, "inputs=8:layout=0_0|w0_0|0_h0|w0_h0|w0+w1_0|w0+w1+w2_0|w0+w1_h0|w0+w1+w2_h0", NULL, filter_graph) < 0) {
        fprintf(stderr, "Could not create xstack_qsv filter\n");
        avfilter_graph_free(&filter_graph);
        return -1;
    }

    // Create format filters
    AVFilterContext *format_ctx1;
    if (avfilter_graph_create_filter(&format_ctx1, format, NULL, "y210le", NULL, filter_graph) < 0) {
        fprintf(stderr, "Could not create format filter y210le\n");
        avfilter_graph_free(&filter_graph);
        return -1;
    }
    AVFilterContext *format_ctx2;
    if (avfilter_graph_create_filter(&format_ctx2, format, NULL, "yuv422p10le", NULL, filter_graph) < 0) {
        fprintf(stderr, "Could not create format filter yuv422p10le\n");
        avfilter_graph_free(&filter_graph);
        return -1;
    }

    // Create buffer sink filter
    if (avfilter_graph_create_filter(buffersink_ctx, buffersink, NULL, NULL, NULL, filter_graph) < 0) {
        fprintf(stderr, "Could not create buffer sink filter\n");
        avfilter_graph_free(&filter_graph);
        return -1;
    }

    // Link filters
    for (int i = 0; i < num_inputs; i++) {
        if (avfilter_link(buffersrc_ctx[i], 0, hwupload_ctx[i], 0) < 0 ||
            avfilter_link(hwupload_ctx[i], 0, scale_ctx[i], 0) < 0 ||
            avfilter_link(scale_ctx[i], 0, xstack_ctx, i) < 0) {
            fprintf(stderr, "Error linking filters for input %d\n", i);
            avfilter_graph_free(&filter_graph);
            return -1;
        }
    }
    if (avfilter_link(xstack_ctx, 0, format_ctx1, 0) < 0 ||
        avfilter_link(format_ctx1, 0, format_ctx2, 0) < 0 ||
        avfilter_link(format_ctx2, 0, *buffersink_ctx, 0) < 0) {
        fprintf(stderr, "Error linking xstack and format filters\n");
        avfilter_graph_free(&filter_graph);
        return -1;
    }

    // Configure the filter graph
    if (avfilter_graph_config(filter_graph, NULL) < 0) {
        fprintf(stderr, "Error configuring the filter graph\n");
        avfilter_graph_free(&filter_graph);
        return -1;
    }

    *graph = filter_graph;
    return 0;
}

// Function to start processing
void start_processing(const char *input_filenames[], int num_inputs, const char *output_filename) {
    AVFormatContext *input_fmt_ctx[num_inputs];
    AVFormatContext *output_fmt_ctx = NULL;
    AVCodecContext *dec_ctx[num_inputs];
    AVCodecContext *enc_ctx = NULL;
    AVFilterGraph *filter_graph = NULL;
    AVFilterContext *buffersrc_ctx[num_inputs];
    AVFilterContext *buffersink_ctx = NULL;
    int ret;

    // Open input files
    for (int i = 0; i < num_inputs; i++) {
        AVDictionary *options = NULL;
        av_dict_set(&options, "p_port", "${VFIO_PORT_PROC}", 0);
        av_dict_set(&options, "p_sip", "192.168.2.2", 0);
        av_dict_set(&options, "p_rx_ip", "192.168.2.1", 0);
        av_dict_set(&options, "udp_port", "20000", 0);
        av_dict_set(&options, "payload_type", "112", 0);
        av_dict_set(&options, "fps", "25", 0);
        av_dict_set(&options, "pix_fmt", "yuv422p10le", 0);
        av_dict_set(&options, "video_size", "1920x1080", 0);

        if (avformat_open_input(&input_fmt_ctx[i], input_filenames[i], av_find_input_format("mtl_st20p"), &options) < 0) {
            fprintf(stderr, "Could not open input file %s\n", input_filenames[i]);
            av_dict_free(&options);
            goto end;
        }
        av_dict_free(&options);

        // Find stream info
        if (avformat_find_stream_info(input_fmt_ctx[i], NULL) < 0) {
            fprintf(stderr, "Could not find stream info for input file %s\n", input_filenames[i]);
            goto end;
        }

        // Find video stream
        int video_stream_index = av_find_best_stream(input_fmt_ctx[i], AVMEDIA_TYPE_VIDEO, -1, -1, NULL, 0);
        if (video_stream_index < 0) {
            fprintf(stderr, "Could not find video stream in input file %s\n", input_filenames[i]);
            goto end;
        }

        // Set up decoder
        AVStream *video_stream = input_fmt_ctx[i]->streams[video_stream_index];
        AVCodec *dec = avcodec_find_decoder(video_stream->codecpar->codec_id);
        if (!dec) {
            fprintf(stderr, "Could not find decoder for input file %s\n", input_filenames[i]);
            goto end;
        }
        dec_ctx[i] = avcodec_alloc_context3(dec);
        if (!dec_ctx[i]) {
            fprintf(stderr, "Could not allocate decoder context for input file %s\n", input_filenames[i]);
            goto end;
        }
        if (avcodec_parameters_to_context(dec_ctx[i], video_stream->codecpar) < 0) {
            fprintf(stderr, "Could not copy codec parameters to context for input file %s\n", input_filenames[i]);
            goto end;
        }
        if (avcodec_open2(dec_ctx[i], dec, NULL) < 0) {
            fprintf(stderr, "Could not open decoder for input file %s\n", input_filenames[i]);
            goto end;
        }
    }

    // Set up filter graph
    if (setup_filter_graph(&filter_graph, buffersrc_ctx, &buffersink_ctx, dec_ctx, num_inputs) < 0) {
        fprintf(stderr, "Could not set up filter graph\n");
        goto end;
    }

    // Set up output format context
    AVDictionary *output_options = NULL;
    av_dict_set(&output_options, "p_port", "${VFIO_PORT_PROC}", 0);
    av_dict_set(&output_options, "p_sip", "192.168.2.2", 0);
    av_dict_set(&output_options, "p_tx_ip", "192.168.2.3", 0);
    av_dict_set(&output_options, "udp_port", "20000", 0);
    av_dict_set(&output_options, "payload_type", "112", 0);
    av_dict_set(&output_options, "fps", "25", 0);
    av_dict_set(&output_options, "pix_fmt", "yuv422p10le", 0);

    avformat_alloc_output_context2(&output_fmt_ctx, av_find_output_format("mtl_st20p"), NULL, output_filename);
    if (!output_fmt_ctx) {
        fprintf(stderr, "Could not create output context\n");
        av_dict_free(&output_options);
        goto end;
    }
    AVOutputFormat *output_fmt = output_fmt_ctx->oformat;

    // Create output stream
    AVStream *out_stream = avformat_new_stream(output_fmt_ctx, NULL);
    if (!out_stream) {
        fprintf(stderr, "Could not create output stream\n");
        av_dict_free(&output_options);
        goto end;
    }

    // Set up encoder
    AVCodec *enc = avcodec_find_encoder(AV_CODEC_ID_H264);
    if (!enc) {
        fprintf(stderr, "Could not find encoder\n");
        av_dict_free(&output_options);
        goto end;
    }
    enc_ctx = avcodec_alloc_context3(enc);
    if (!enc_ctx) {
        fprintf(stderr, "Could not allocate encoder context\n");
        av_dict_free(&output_options);
        goto end;
    }
    enc_ctx->height = dec_ctx[0]->height / 2; // Because of scale_qsv
    enc_ctx->width = dec_ctx[0]->width / 4 * 4; // Because of scale_qsv and xstack
    enc_ctx->sample_aspect_ratio = dec_ctx[0]->sample_aspect_ratio;
    enc_ctx->pix_fmt = AV_PIX_FMT_YUV422P10LE;
    enc_ctx->time_base = (AVRational){1, 25};
    if (output_fmt->flags & AVFMT_GLOBALHEADER) {
        enc_ctx->flags |= AV_CODEC_FLAG_GLOBAL_HEADER;
    }
    if (avcodec_open2(enc_ctx, enc, NULL) < 0) {
        fprintf(stderr, "Could not open encoder\n");
        av_dict_free(&output_options);
        goto end;
    }
    if (avcodec_parameters_from_context(out_stream->codecpar, enc_ctx) < 0) {
        fprintf(stderr, "Could not copy encoder parameters to stream\n");
        av_dict_free(&output_options);
        goto end;
    }

    // Open output file
    if (!(output_fmt->flags & AVFMT_NOFILE)) {
        if (avio_open2(&output_fmt_ctx->pb, output_filename, AVIO_FLAG_WRITE, NULL, &output_options) < 0) {
            fprintf(stderr, "Could not open output file\n");
            av_dict_free(&output_options);
            goto end;
        }
    }
    av_dict_free(&output_options);

    // Write output file header
    if (avformat_write_header(output_fmt_ctx, NULL) < 0) {
        fprintf(stderr, "Could not write output file header\n");
        goto end;
    }

    // Allocate frames and packets
    AVFrame *frame = av_frame_alloc();
    AVFrame *filt_frame = av_frame_alloc();
    AVPacket *pkt = av_packet_alloc();

    // Read, decode, filter, encode, and write frames
    for (int i = 0; i < num_inputs; i++) {
        while (av_read_frame(input_fmt_ctx[i], pkt) >= 0) {
            if (pkt->stream_index == 0) {
                ret = avcodec_send_packet(dec_ctx[i], pkt);
                if (ret < 0) {
                    fprintf(stderr, "Error sending packet to decoder for input file %s\n", input_filenames[i]);
                    break;
                }
                while (ret >= 0) {
                    ret = avcodec_receive_frame(dec_ctx[i], frame);
                    if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) {
                        break;
                    } else if (ret < 0) {
                        fprintf(stderr, "Error receiving frame from decoder for input file %s\n", input_filenames[i]);
                        goto end;
                    }

                    // Send frame to filter graph
                    if (av_buffersrc_add_frame(buffersrc_ctx[i], frame) < 0) {
                        fprintf(stderr, "Error sending frame to filter graph for input file %s\n", input_filenames[i]);
                        goto end;
                    }

                    // Receive filtered frame
                    while (1) {
                        ret = av_buffersink_get_frame(buffersink_ctx, filt_frame);
                        if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) {
                            break;
                        } else if (ret < 0) {
                            fprintf(stderr, "Error receiving frame from filter graph\n");
                            goto end;
                        }

                        // Encode filtered frame
                        ret = avcodec_send_frame(enc_ctx, filt_frame);
                        if (ret < 0) {
                            fprintf(stderr, "Error sending frame to encoder\n");
                            goto end;
                        }
                        while (ret >= 0) {
                            ret = avcodec_receive_packet(enc_ctx, pkt);
                            if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) {
                                break;
                            } else if (ret < 0) {
                                fprintf(stderr, "Error receiving packet from encoder\n");
                                goto end;
                            }

                            // Write packet to output file
                            pkt->stream_index = out_stream->index;
                            av_packet_rescale_ts(pkt, enc_ctx->time_base, out_stream->time_base);
                            if (av_interleaved_write_frame(output_fmt_ctx, pkt) < 0) {
                                fprintf(stderr, "Error writing packet to output file\n");
                                goto end;
                            }
                            av_packet_unref(pkt);
                        }
                        av_frame_unref(filt_frame);
                    }
                    av_frame_unref(frame);
                }
            }
            av_packet_unref(pkt);
        }
    }

    // Write output file trailer
    av_write_trailer(output_fmt_ctx);

end:
    // Cleanup
    av_frame_free(&frame);
    av_frame_free(&filt_frame);
    av_packet_free(&pkt);
    for (int i = 0; i < num_inputs; i++) {
        avcodec_free_context(&dec_ctx[i]);
        avformat_close_input(&input_fmt_ctx[i]);
    }
    avcodec_free_context(&enc_ctx);
    if (output_fmt_ctx && !(output_fmt_ctx->oformat->flags & AVFMT_NOFILE)) {
        avio_closep(&output_fmt_ctx->pb);
    }
    avformat_free_context(output_fmt_ctx);
    if (filter_graph) {
        avfilter_graph_free(&filter_graph);
    }
}

int main() {
    av_register_all();
    avfilter_register_all();
    
    const char *input_filenames[] = {
        "[0]",
        "[1]",
        "[2]",
        "[3]",
        "[4]",
        "[5]",
        "[6]",
        "[7]"
    };
    int num_inputs = sizeof(input_filenames) / sizeof(input_filenames[0]);

    // Start processing
    start_processing(input_filenames, num_inputs, "-");

    return 0;
}