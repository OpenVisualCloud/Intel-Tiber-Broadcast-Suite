#include <iostream>

#include "struct_new.hpp"

#ifndef _FFMPEG_PIPELINE_GENERATOR_H_
#define _FFMPEG_PIPELINE_GENERATOR_H_

/**
 * @brief Generates an FFmpeg pipeline string based on the provided configuration.
 *
 * This function takes a configuration object and constructs an FFmpeg pipeline string
 * that can be used to execute FFmpeg commands. The configuration object contains various
 * settings and parameters that dictate how the pipeline should be constructed.
 *
 * @param config The configuration object containing settings for the FFmpeg pipeline.
 * @param pipeline_string The string where the generated FFmpeg pipeline will be stored.
 * @return int Returns 0 on success, or a non-zero error code on failure.
 */
int ffmpeg_generate_pipeline(Config &config, std::string &pipeline_string);

#endif // _FFMPEG_PIPELINE_GENERATOR_H_