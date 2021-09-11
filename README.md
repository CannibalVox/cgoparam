# cgoparam
Fast, thread-safe arena allocators

### What is CgoParam?

CgoParam is a refinement & specialization of [cgoalloc](https://github.com/CannibalVox/cgoalloc) - a general-purpose allocation proxying library for cgo.

Cgoalloc offers paging fixed-buffer allocators that can reduce the dependence on C.malloc and C.free (as well as avoiding expensive cgocheckpointers by only passing C memory to C).  It also allows various allocator types to be composed together to create interesting allocation & retainment strategies.

However, cgoalloc requires allocators to be composed from several parts to be effective, and it's not quite as fast as it could be.  The FixedBufferAllocator can perform an allocation and free in 10ns, but an arena allocator built on top of it can take up to 30ns.  That's for incredibly small allocations, too, and it's slower when the arena allocator is built on top of a 3-part allocator to ensure it can handle memory of any size.  cgoalloc structures are also not thread-safe.

As a result of these limitations, cgoalloc can support any allocation and retainment strategy you might like.

By contrast, cgoparam is built for a single purpose- temporary allocations of cgo parameter pointers, assigned just before a cgo call and freed just after.  It uses a sync.Pool in order to make the library thread-safe (allocators are not thread-safe, but each thread can freely pull and return allocators), meaning that cgoparam can be an implementation detail of your cgo wrapper library- a thing your users don't have to worry about.

For this role, cgo's performance is perfect:

#### Small Batch Test

Four allocations are assigned and then freed: 1 byte, 10 bytes, 100 bytes, and 1000 bytes

```
BenchmarkRawCgoSmallBatch-16             3185971               374.9 ns/op
BenchmarkCgoparamSmallBatch-16          61446448                19.74 ns/op
```

#### Multi-Page Test

Ten allocations of 1000 bytes are made and then freed

```
BenchmarkRawCgoMultipage-16               945380              1150 ns/op
BenchmarkCgoparamMultipage-16            4181353               293.5 ns/op
```

Pages are 4096 bytes- any additional pages allocated during the life of a single allocator must both be malloc'd and freed, but if you stay within the 4096 byte limit, your performance will more closely resemble the first example: 5ns per alloc/free!

### Example

```go
    allocator := cgoparam.GetAlloc()
    defer cgoparam.ReturnAlloc()
    
    createInfo := (*C.VkBufferCreateInfo)(allocator.Malloc(C.sizeof_struct_VkBufferCreateInfo))
    createInfo.sType = C.VK_STRUCTURE_TYPE_BUFFER_CREATE_INFO
    createInfo.flags = 0
    createInfo.size = C.VkDeviceSize(o.BufferSize)
    createInfo.usage = C.VkBufferUsageFlags(o.Usages)
    createInfo.sharingMode = C.VkSharingMode(o.SharingMode)
    
    queueFamilyCount := len(o.QueueFamilyIndices)
    createInfo.queueFamilyIndexCount = C.uint32_t(queueFamilyCount)
    createInfo.pQueueFamilyIndices = nil
    
    if queueFamilyCount > 0 {
        indicesPtr := (*C.uint32_t)(allocator.Malloc(queueFamilyCount * int(unsafe.Sizeof(C.uint32_t(0)))))
        indicesSlice := ([]C.uint32_t)(unsafe.Slice(indicesPtr, queueFamilyCount))
        
        for i := 0; i < queueFamilyCount; i++ {
            indicesSlice[i] = C.uint32_t(o.QueueFamilyIndices[i])
        }
        
        createInfo.pQueueFamilyIndices = indicesPtr
    }
    
    return unsafe.Pointer(createInfo), nil
}
```
