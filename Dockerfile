FROM golang:1.13.4-alpine

MAINTAINER Cache Lab <hello@cachelab.co>

COPY aws_cost_waste /bin/aws_cost_waste

USER nobody

ENTRYPOINT ["/bin/aws_cost_waste"]
