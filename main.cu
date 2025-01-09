#include <iostream>
#include <vector>
#include <algorithm>       // For std::min
#include "cuda_runtime.h"
#include "device_launch_parameters.h"

// If you use nlohmann::json anywhere, make sure to include it properly:
// #include "json.hpp"
// using json = nlohmann::json;

using namespace std;

#define CHUNK_SIZE 10
#define THREADS_PER_BLOCK 10

/**
 * Kernel to find the maximum value in a chunk of `numbers`.
 * Each thread compares its element and does an atomicMax with the shared device max pointer.
 */
__global__ void processMaxNumberKernel(const int* numbers, int size, int* d_max)
{
    int idx = blockDim.x * blockIdx.x + threadIdx.x;
    if (idx < size)
    {
        // Use atomicMax to update the current maximum in a thread-safe manner.
        atomicMax(d_max, numbers[idx]);
    }
}

/**
 * Allocate and fill a vector of size `count` with sequential integers [0, 1, 2, ..., count-1].
 */
vector<int> AllocateSquaredNumbers(int count)
{
    vector<int> results;
    results.reserve(count);

    for (int i = 0; i < count; i++)
    {
        results.push_back(i * i);
    }

    return results;
}

int main()
{
    int count = 100;
    int maxNumber = 0;  // Host-side running maximum

    // Allocate numbers [0..99] on the host
    vector<int> numbers = AllocateSquaredNumbers(count);

    for (int i = 0; i < count; i += CHUNK_SIZE)
    {
        int currentChunkSize = min(CHUNK_SIZE, count - i);

        // Allocate device memory for the current chunk
        int* d_results;
        cudaMalloc(&d_results, currentChunkSize * sizeof(int));

        // Copy chunk from host to device
        cudaMemcpy(d_results, &numbers[i],
            currentChunkSize * sizeof(int),
            cudaMemcpyHostToDevice);

        // Allocate device memory for the chunk maximum and initialize it with the current host max
        int* d_chunkMax;
        cudaMalloc(&d_chunkMax, sizeof(int));
        cudaMemcpy(d_chunkMax, &maxNumber, sizeof(int), cudaMemcpyHostToDevice);

        // Number of blocks needed
        int blocks = (currentChunkSize + THREADS_PER_BLOCK - 1) / THREADS_PER_BLOCK;

        // Launch kernel to update d_chunkMax via atomicMax
        processMaxNumberKernel << <blocks, THREADS_PER_BLOCK >> > (d_results, currentChunkSize, d_chunkMax);
        cudaDeviceSynchronize();

        // Read back the chunk's max from device
        int chunkMax = 0;
        cudaMemcpy(&chunkMax, d_chunkMax, sizeof(int), cudaMemcpyDeviceToHost);

        // Update global max if needed
        if (chunkMax > maxNumber)
        {
            maxNumber = chunkMax;
        }

        // Clean up device memory
        cudaFree(d_results);
        cudaFree(d_chunkMax);
    }

    cout << "Max number: " << maxNumber << endl;

    return 0;
}