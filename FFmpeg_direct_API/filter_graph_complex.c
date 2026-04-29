#include <libavfilter/avfilter.h>

#include "filter_graph_complex.h"

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
