FROM alpine:3.11
ADD omo.msa.acm /usr/bin/omo.msa.acm
ENV MSA_REGISTRY_PLUGIN
ENV MSA_REGISTRY_ADDRESS
ENTRYPOINT [ "omo.msa.acm" ]
