#!/bin/bash

export RANDOM_GUESS_COUNT="1"
export MAX_GUESS_COUNT="9"
export GENERATION_MIN="1"
export GENERATION_MAX="9"
# export HIDE_HEIGHT_COLUMN=""
# export HIDE_WEIGHT_COLUMN=""
# export HIDE_GENERATION_COLUMN=""
# export HIDE_REMAINING_TYPES=""
# export HIDE_SECOND_TYPE_COLUMN=""

go build
./chonkle