rm -rf CMakeFiles/
rm -rf cmake_install.cmake
rm -rf CMakeCache.txt
rm -rf Makefile

cmake -G "Unix Makefiles" -DANDROID_ABI=armeabi-v7a -DANDROID_PLATFORM=android-21 -DCMAKE_BUILD_TYPE=Release -DANDROID_NDK=${ANDROID_NDK_HOME} -DCMAKE_TOOLCHAIN_FILE=${ANDROID_NDK_HOME}/build/cmake/android.toolchain.cmake -DANDROID_NATIVE_API_LEVEL=21 -DANDROID_STL=c++_shared -DANDROID_TOOLCHAIN=clang .
make clean
make -j4

cmake -G "Unix Makefiles" -DANDROID_ABI=arm64-v8a -DANDROID_PLATFORM=android-21 -DCMAKE_BUILD_TYPE=Debug -DANDROID_NDK=${ANDROID_NDK_HOME} -DCMAKE_TOOLCHAIN_FILE=${ANDROID_NDK_HOME}/build/cmake/android.toolchain.cmake -DANDROID_NATIVE_API_LEVEL=21 -DANDROID_STL=c++_shared -DANDROID_TOOLCHAIN=clang .
make clean
make -j4

cmake -G "Unix Makefiles" -DANDROID_ABI=x86 -DANDROID_PLATFORM=android-21 -DCMAKE_BUILD_TYPE=Release -DANDROID_NDK=${ANDROID_NDK_HOME} -DCMAKE_TOOLCHAIN_FILE=${ANDROID_NDK_HOME}/build/cmake/android.toolchain.cmake -DANDROID_NATIVE_API_LEVEL=21 -DANDROID_STL=c++_shared -DANDROID_TOOLCHAIN=clang .
make clean
make -j4

cmake -G "Unix Makefiles" -DANDROID_ABI=x86_64 -DANDROID_PLATFORM=android-21 -DCMAKE_BUILD_TYPE=Release -DANDROID_NDK=${ANDROID_NDK_HOME} -DCMAKE_TOOLCHAIN_FILE=${ANDROID_NDK_HOME}/build/cmake/android.toolchain.cmake -DANDROID_NATIVE_API_LEVEL=21 -DANDROID_STL=c++_shared -DANDROID_TOOLCHAIN=clang .
make clean
make -j4