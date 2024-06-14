ARG IMAGE_TAG_NAME
FROM ${IMAGE_TAG_NAME}_build_stage
SHELL ["/bin/bash", "-e", "-o", "pipefail", "-c"]
ENV DEBIAN_FRONTEND=noninteractive


ARG ABI_VER=2.5.0
# Install Coverity
RUN wget -nv -P /tmp https://ubit-artifactory-or.intel.com/artifactory/coverity-or-local/Archive/2023.3.0/cov-analysis-linux64-2023.3.0.sh \
    && wget -nv -P /tmp https://ubit-artifactory-or.intel.com/artifactory/coverity-or-local/Enterprise/license.dat
RUN chmod +x /tmp/cov-analysis-linux64-2023.3.0.sh \
    && /tmp/cov-analysis-linux64-2023.3.0.sh \
    --installation.dir=/OWR/Tools/coverity/2023.3.0/ \
    --license.agreement=agree \
    --license.cov.path=/tmp/license.dat \
    --license.region=0 \
    --license.type.choice=0 \
    --component.sdk=false \
    # NOTE: Be careful, if any characters are changed below you risk the nonsensical error:
    #   The boolean argument --component.skip.documentation had value true​.
    #   The value should be either true or false.
    --component.skip.documentation=true -q \
    && rm -f /tmp/cov-analysis-linux64-2023.3.0.sh

# Install Coverity Reports
RUN wget -nv -P /tmp https://ubit-artifactory-or.intel.com/artifactory/coverity-or-local/Archive/2023.3.0/reports/cov-reports-linux64-2023.3.0.sh
RUN chmod +x /tmp/cov-reports-linux64-2023.3.0.sh \
    && /tmp/cov-reports-linux64-2023.3.0.sh \
    -q -console -dir /OWR/Tools/cov_report/2023.3.0
ENV PATH "${PATH}:/OWR/Tools/cov_report/2023.3.0/bin"

# Install ABI
RUN apt-get update -y \
 && apt-get upgrade -y \
 && apt-get install --fix-missing --no-install-recommends -y \
    python3 \
    python3-pip \
    python3-venv \
    python3-dev \
    libkrb5-dev \
    wget \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*
RUN python3 -m pip install --upgrade pip setuptools wheel \
 && python3 -m pip install abi==${ABI_VER} --extra-index-url=https://ubit-artifactory-or.intel.com/artifactory/api/pypi/one-windows-pypi-local/simple

COPY opt/* /usr/bin
WORKDIR /opt/
RUN chmod +x /usr/bin/coverity && chmod +x /usr/bin/build.sh

CMD [ "coverity "]
