# Understanding

为了实现一个性能超过`sort.Slice()`函数的归并排序算法，我们需要利用多核CPU，实现一个并行归并排序。对于每个并行的部分，我们可以直接利用`sort.Slice()`进行排序。对于归并的部分，我们需要实现一个归并算法。由于不同的机器的CPU核心数可能不同，我们需要实现一个多路归并排序算法。可参考[Wikipedia的介绍][k-way-merge]进行实现。

# Design

要编写一个充分利用每一个CPU核心的并行排序算法，我们需要让每个goroutine并行地运行在每一个操作系统线程上。因此，需要创建数量为`runtime.NumCPU()`的goroutine进行分段并行排序。最后，通过同步操作在每个“排序goroutine”运行结束后，将所有的部分有序的数组进行多路归并，得出最后完整排序好的数组。

# Build

为了充分利用多核CPU，`kway.Sort()`函数创建了`n`个goroutine，其中`n`等于`runtime.NumCPU()`，也就是CPU核心的个数。通过查看trace的goroutine analysis，我们可以看到当`n`等于`runtime.NumCPU()`时，每个goroutine的运行时间几乎全部都是在运行（Execution）状态。
![goroutine_analysis_v0](./_pprof/img/goroutine_analysis_v0.png)

当我们修改代码使`n`等于两倍的`runtime.NumCPU()`时，
```diff
diff --git a/tidb/mergesort/kway/merge.go b/tidb/mergesort/kway/merge.go
index f173e9b..fa57183 100644
--- a/tidb/mergesort/kway/merge.go
+++ b/tidb/mergesort/kway/merge.go
@@ -10,7 +10,7 @@ import (
 // Sort sorts the array with multiple goroutines using K-way Merge Sort algorithm.
 func Sort(data []int64) {
        var (
-               n = runtime.NumCPU()
+               n = 2 * runtime.NumCPU()
        )

        if len(data) < 10*n {
```

每个goroutine的运行时间会有几乎一半用于等待调度（Schedular wait）状态。
![goroutine_analysis_v0_double_goroutine](./_pprof/img/goroutine_analysis_v0_double_goroutine.png)

在实现goroutine同步部分，`kway.Sort()`使用`sync.WaitGroup`来在排序goroutine与多路归并goroutine之间进行同步。

在实现多路归并算法部分，`kway.Sort()`使用[Tournament Tree算法][tournament]。Tournamenet Tree算法中，维护一个loser tree，其底层使用一维数组实现。树的存储方式与堆的底层存储方式类似，并在不需要存储节点的`tree[0]`上用于存储每次的winner。*TODO: 重新comment算法部分后，再refine一下。*

# Optimizing

第一次实现，`make bench`结果已经比`sort.Slice()`好。

具体分析主要瓶颈在`sort.Slice()`上，*TODO：真的吗？查一下如何分析瓶颈*

*TODO: 整理算法代码后，预估性能再看profile结果*

[k-way-merge]: https://en.wikipedia.org/wiki/K-way_merge_algorithm
[tournament]: https://en.wikipedia.org/wiki/K-way_merge_algorithm#Tournament_Tree
