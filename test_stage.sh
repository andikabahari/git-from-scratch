#!/bin/sh

case $1 in
  '1') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='init'
    ;;
  '2') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='read_blob'
    ;;
  '3') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='create_blob'
    ;;
  '4') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='read_tree'
    ;;
  '5') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='write_tree'
    ;;
  '6') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='create_commit'
    ;;
  '7') 
    CODECRAFTERS_CURRENT_STAGE_SLUG='clone_repository'
    ;;
  *)
    echo 'Invalid stage'
    exit
    ;;
esac

cd git-tester
go build -o ../git-go/test.out ./cmd/tester

cd ../git-go
CODECRAFTERS_SUBMISSION_DIR=$(pwd) \
CODECRAFTERS_CURRENT_STAGE_SLUG=${CODECRAFTERS_CURRENT_STAGE_SLUG} \
./test.out
rm ./test.out