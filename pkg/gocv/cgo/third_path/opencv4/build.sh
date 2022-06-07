#!/bin/bash

rm -rf bin include lib64 share turbo build
rm -rf libjpeg-turbo-2.1.1 opencv-4.5.3
tar -zxvf libjpeg-turbo-2.1.1.tar.gz
tar -zxvf opencv-4.5.3.tar.gz
sed -i 'N;2 a add_compile_options(-fPIC)' libjpeg-turbo-2.1.1/CMakeLists.txt

PREFIX=`pwd`
mkdir build; cd build
cmake ../libjpeg-turbo-2.1.1 -DCMAKE_INSTALL_PREFIX=$PREFIX/turbo -DCMAKE_BUILD_TYPE=RELEASE
make -j8 && make install
rm -rf *

cmake ../opencv-4.5.3 -DCMAKE_INSTALL_PREFIX=$PREFIX -DCMAKE_BUILD_TYPE=RELEASE -DBUILD_SHARED_LIBS=ON -DBUILD_PNG=ON -DBUILD_JASPER=ON -DBUILD_TIFF=ON -DWITH_FFMPEG=OFF -DBUILD_PERF_TESTS=OFF -DBUILD_TESTS=OFF -DBUILD_JPEG=OFF -DBUILD_OPENEXR=ON -DWITH_CUDA=OFF -DOPENCV_GENERATE_PKGCONFIG=ON -DWITH_JPEG=ON -DJPEG_INCLUDE_DIR=$PREFIX/turbo/include -DJPEG_LIBRARY=$PREFIX/turbo/lib64/libjpeg.a
make -j8 && make install

cd ..
cp -rf include/opencv4/* include/
rm -rf include/opencv4
