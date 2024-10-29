ARG IMAGE_TAG_NAME
FROM ${IMAGE_TAG_NAME}
SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c"]
ENV DEBIAN_FRONTEND=noninteractive

ARG COVERITY_VER=2023.3.4
ARG ABI_VER=3.1.0
# Install Coverity
RUN wget -nv -P /tmp https://ubit-artifactory-or.intel.com/artifactory/coverity-or-local/Archive/${COVERITY_VER}/cov-analysis-linux64-${COVERITY_VER}.sh \
    && wget -nv -P /tmp https://ubit-artifactory-or.intel.com/artifactory/coverity-or-local/Enterprise/license.dat
RUN chmod +x /tmp/cov-analysis-linux64-${COVERITY_VER}.sh \
    && /tmp/cov-analysis-linux64-${COVERITY_VER}.sh \
    --installation.dir=/OWR/Tools/coverity/${COVERITY_VER}/ \
    --license.agreement=agree \
    --license.cov.path=/tmp/license.dat \
    --license.region=0 \
    --license.type.choice=0 \
    --component.sdk=false \
    # NOTE: Be careful, if any characters are changed below you risk the nonsensical error:
    #   The boolean argument --component.skip.documentation had value true​.
    #   The value should be either true or false.
    --component.skip.documentation=true -q \
    && rm -f /tmp/cov-analysis-linux64-${COVERITY_VER}.sh

# Install Coverity Reports
RUN wget -nv -P /tmp https://ubit-artifactory-or.intel.com/artifactory/coverity-or-local/Archive/${COVERITY_VER}/reports/cov-reports-linux64-${COVERITY_VER}.sh
RUN chmod +x /tmp/cov-reports-linux64-${COVERITY_VER}.sh \
    && /tmp/cov-reports-linux64-${COVERITY_VER}.sh \
    -q -console -dir /OWR/Tools/cov_report/${COVERITY_VER}
ENV PATH "${PATH}:/OWR/Tools/cov_report/${COVERITY_VER}/bin:/root/.local/bin"

# Install ABI
RUN apt-get update -y \
 && apt-get upgrade -y \
 && apt-get install --fix-missing --no-install-recommends -y \
    python3 \
    python3-pip \
    python3-venv \
    python3-dev \
    libkrb5-dev \
    libfreetype6-dev \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*
RUN python3 -m pip install --break-system-packages --upgrade setuptools \
 && python3 -m pip install --break-system-packages --user abi==${ABI_VER} --extra-index-url=https://ubit-artifactory-or.intel.com/artifactory/api/pypi/one-windows-pypi-local/simple

COPY opt/* /usr/bin
WORKDIR /tmp/host
RUN chmod +x /usr/bin/coverity && \
    chmod +rw /usr/bin/report_cfg.yaml && \
    chmod +x /usr/bin/build_ffmpeg.sh && \
    chmod +x /usr/bin/build_vsr.sh && \
    chmod +x /usr/bin/build_onevpl.sh && \
    chmod +x /usr/bin/build_mtl.sh

CMD ["/bin/bash"]
