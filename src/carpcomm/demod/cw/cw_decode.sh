#!/bin/bash

python src/carpcomm/demod/cw_filter.py $1 \
    | ./bin/demod --satellite_id=$2