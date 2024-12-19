#ifndef _FILTER_GRAPH_COMPLEX_H_
#define _FILTER_GRAPH_COMPLEX_H_

#include <libavfilter/avfilter.h>
#include <libavfilter/buffersink.h>
#include <libavfilter/buffersrc.h>

// Function to set up the filter graph
int setup_filter_graph(AVFilterGraph **graph, AVFilterContext **buffersrc_ctx,
                       AVFilterContext **buffersink_ctx, AVCodecContext *dec_ctx[],
                       int num_inputs);

#endif _FILTER_GRAPH_COMPLEX_H_
