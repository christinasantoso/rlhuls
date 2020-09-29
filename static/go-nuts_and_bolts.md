[**Tags**] go; database; kv; boltdb
[**title**] Go-nuts and Bolts: An Introduction to BoltDB
[**speaker**] Tommi Virtanen
[**date**] 2015-02-18


[**transcript**]
So Go-nuts and Bolts, Go-nuts is the remaining list and Bolts is this awesome database library, key-value store library. It’s now GitHub-- there’s an interesting delay with my slides, if that (slides) doesn’t update, tell me, it updates on my screen. So Bolt is on GitHub as boltdb/bolt and it’s a project primarily by Ben Johnson who was an awesome dude, he writes really neat code. If you’re looking to understand something that’s a little fiddly with a low-level detail but really prettily written, that source code is really a good experience to go through. Oh, my involvement with bolt has been mostly to be a sounding board for API design and I also snuck in a couple of features. I’m a heavy user, and I’m here as a user telling you that bolt is awesome. Bolt is a key-value store; it's not a full database, it’s not a sequel database, definitely not. It’s a key-value store, this key has this value.


[**slide**] 1:30
```
key-value store
```
db, err := bolt.Open(path, 0644, nil)
if err != nil {
        log.Fatal(err)
}
defer db.Close()
```
```
if err := bucket.Put([]byte(“answer”), []byte(“hello”)); err !=nil {
        return err
}
```
```
val := bucket.Get([]byte(“answer”))
if val == nil {
        // not found
        return errors.New(“no answer”)
}
fmt.Println(val)
```
all [ ]byte
we’ll come back to bucket
```


[**transcript**]
Basically, you open a file, you do put operation where this key is this value. You do get operation where you’re giving the key and all these are bite-slices, so there’s no like numerical data types, none of that stuff that belongs in your application or a wrapper you built on top of all that actually understands your data types. As far as all this concern everything is bite-sliced and we’ll talk more about buckets later.






[**slide**] 2:01
```
Bolt
```
key-value store
pure Go
```


[**transcript**]
So bolt is pure go and that’s one reason why it’s so amazing, like why I enjoy using it.


[**slide**] 2:09
pure Go
```
inspired by LMDB
simplicity as a virtue
LMDB is performance over just about everything, has a lot more features
see github.com/szferi/gomdb for an LMDB CGo wrapper
```


[**transcript**]
It was inspired by a project called LMDB which was written for the OpenGL tab so there’s a bucket for that product and it has the same basic design but it is not a cohort/Go port of LMDB. It has different design goals and they’re starting to deviate from each other. Bolt is explicitly as simple as it possibly can be, even if that costs a little bit of performance, bolt chooses to be simple and understandable. LMDB is performance nuts written in C or C++-- I forget which one- and explicitly meant for that sort of stuff. Bolt is a simple pure go and a joy to use. Bolt is a library, it’s not on network service, it’s not something that runs independently of your application. it’s part of your application.


[**slide**] 3:01
```
library
```
not something you talk to over the network
no client and server parts
stored in a local file
you could write the network bits, if you wanted;
usually better to provide more value than just put/get
```






[**transcript**]
And it is actually a file format that you use through this library, so it’s a single file which makes it really easy to manage and you can write network based on top of it but it’s probably better to provide more value than just put/get operation and influencerDB which is the second talk of today is kind of example using Bolt in your own application. 


[**slide**] 3:27
```
Bolt
```
key-value store
pure Go
library
single process…
```


[**transcript**]
Bolt databases are meant to be used by a single process


[**slide**] 3:32
```
single process
```
technically, single owner
one bolt.Open at a time
```


[**transcript**]
So there’s a single owner of the database one at a time, so if somebody has the database open, nobody else can open it at the same time


[**slide**] 3:43
```
Bolt
```
key-value store
pure Go
library
single process
ordered ...
```


[**transcript**]
And the key-value store is ordered so you can work the keys in order


[**slide**] 3:50
```
ordered key-value store
```
use keys that sort the right way
```
```
// makeKey is a toy key encoder. Use something like encoding/binary.
// BigEndian.PutUint32 in real code.
func makeKey(n uint32) []byte {
        return []byte(fmt.Sprintf(“%o4d”, n))
}
```
protip: don’t rely on size of int
```
        for i := 0 i < 100; i++ {
                key := makeKey(uint32(rand.Intn(10000)))
                if err := bucket.Put(key, []byte(‘’)); err != nil {
                }
        }
```
        C := bucket.Cursor()
        for k, v := c.First(); k != nil; k, v = c.Next() {
                fmt.Printf(“%q %v\n”, k, v)
        }
```


[**transcript**]
So here’s an example of-- the first thing you have to remember is the keys and the values are just bit-slices, So to be able to have a meaningful order we need to encode our values in a way that preserves the order. So here we’re doing zero-padded decimal strings as the key. In your real application, you probably do something a little bit more binary, but the same idea is there, so for example big-endian is an order-preserving for integers, [inaudible4:31] integers, So we go through and put in a hundred random keys from 0 to 9999 and then we walk through them in a form and what you get out is you get a listing off. So then we get the main order.




[**slide**] 5:04
```
ordered key-value store: prefix
```
        prefix := []byte(“40”)
        c := bucket.Cursor()
        // while not at the end, and still has wanted prefix
        for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v =c.Next() {
        }
```






[**transcript**]
So an ordered key-value store also lets you access data by others— like gimmicks. Like you can ask for anything with this prefix, so we stored 4 digits decimal numbers and we can ask for anything that begins with a four and a zero. In our 100 random numbers, there are two of those(?)


[**slide**] 5:27
```
ordered key-value store: range
```
        var min uint32 = 4242
        var max uint32 = 4730
maxKey := makeKey(max)
c := bucket.Cursor()
// while not at the end, and not past max
for K, v := c.Seek(makeKey(min)); k != nil && bytes.Compare(k, maxKey) <= 0; k, v = c.Next [Cropped]
fmt.Printf(“%q %v\n”, k, v)
}
```
makeKey zero-padding matters here; order of marshaled keys must match the desired order
```
 
[**transcript**]
You can ask for a range by encoding both the beginning and the ending key, the way you encode keys and walking out cursor from the beginning key to the ending key as long as you don’t go past the maximum key. So it’s kind of very simple, but you can see how this is not like running a sequel query, you’re actually like programming it with a very low-level mechanism. This actually feels like I’m using it in-memory data structure like a container of some sort would have this API now just a persistent data structure.


[**slide**] 6:05
```
Bolt
```
key-value store
pure Go
library
single process
ordered 
b+tree…
```


[**transcript**]
So bolt actually is a B+ tree underneath that was stored in the file and because saying B+tree is annoying and technically B three is something slightly different. I just-- what-- I’m gonna say B tree and it’s gonna be a decimal matrix.


[**slide**] 6:20
```
Copy-on-write B tree
* two metadata pages: more recent wins
* copy-on-write data pages
data pages contain
* buckets
* values
* freelist
write transaction commit first syncs dirty pages to disk,
then syncs updated meta page.
read transactions pin pages that would otherwise recycled.
```


[**transcript**]
So it’s a copy-on-write B tree that begins with two metadata pages which one’s newer is the one that wins. So when we’re opening the database recovering from a crash or whatever which one of those two metadata pages has a higher stamp-- generation power on it, that one wins and that page contains pages-- the rest of the stuff. The rest of the stuff is data pages that contain the buckets that we talked about. They have the key system, they have the values and then there’s a list of which pages are free that is also stored on the pages themselves which is a new-trick Think about that for a moment. It keeps track of what's free [inaudible7:00], that sounds like [inaudible], these sound like in-memory-data-structure. Anyway, when you do a transaction-- an update transaction-- a write transaction, first you writes out all the dirty data pages to new locations so that’s the copy-on-write operation. And then it writes over the other-- the redundant better data pages which now become the new one. So every single transaction that you do is to this [inaudible 7:40] to this syncs. Read the transaction pin the data that would otherwise be recycled that holds it from being recycled. But it’s a big copy on write B tree that uses the file as its storage.


[**slide**] 7:53
```
B tree trade-offs vs LSM
compare: Log-Structured Merge-tree (LSM), for example LevelDB
* LSM appends to a log, keeps log content also in memory
* once log is big enough, sorts and dumps into immutable file
* deletions are puts of special tombstone values
* periodically compact immutable files into fewer, bigger ones
* gets often need to look at multiple files
* all disk writes are sequential; designed for HDDs with poor random I/O
* still benefits from SSDs, but ratio is much smaller
best fit:
LSM: write-heavy workloads on spinning discs, latency spikes during compaction ok
Bolt: read-heavy workloads, SSDs
```






[**transcript**]
So it has all the good and bad sides of B tree in general and one of the sorts of copy-on-write alternatives that gets mentioned is LSM like LevelDB for example. And these two are completely different in design but used for similar purposes which kinda lead you to question which one should we use. Well, the basic guideline is that LevelDB was designed to take random writes and convert them into sequential writes as they go to disk. And the thing that benefits from that is spinning discs. So LevelDB is really intended for write-heavy operations running on spinning discs. It benefits from SSDs but it really is designed for spinning disk. B trees tend to have a little of random accesses which while they still hurt with SSD they hurt significantly less. So bolt is really good when running on SSDs and in general the design is flipped all the way up the other way around. So whereas a log-structured merge-tree like LevelDB is write-heavy and read direction slower than writes. Bolt is the exact opposite, as heavy of a great load you can give, it’s just that what is it built for and writes are a little bit more expansive. And the other thing is that Log-Structured merge-tree tends to have this spikes of heavy I/O compaction, where it cleans up the database, whereas a B tree will do that to all the file all the time. So once again it’s a little bit slower from the average offering but it has steady performance not spiking performance every now and then.


[**slide**] 9:44
```
Bolt
```
key-value store
pure Go
library
single process
ordered 
b+tree
single writer...
```
[**transcript**]
We don’t have much time to talk about that now but find me later in this conference to ask me questions.Cause that’s kinda off topic in itself that's just a comparison against another thing that is in a similar role. So one of the reasons that make bolt pleasant to use is because it’s simple and one of the major ways it is simple is because it has a single writer at a time, touching the database. So if you think about your typical-- like count-- you read the count, you increment the number, you store it back. Now, something gonna  – “aha you can’t do that, you would have conflicts with two updates happening at the same time”. Well, Bolt is explicitly designed so that never happens. So, doing the dumbest possible thing, writing the simplest possible code is actually what is meant to happen with Bolt. You don’t need to have this elaborate, for example in the sequel you need to write a query that does the thing you want instead of reading the value into your application, doing something for it, and putting it back in






[**slide**] 10:55
```
single writer 
if function f is updating the database, no one else is
serializable
very easy to reason about
safe to do read-modify-write
limits performance, but that may not matter
```


[**transcript**]
So Bolt is-- it feels more like programming than writing in something yourself, did that make sense? and I actually prefer that because my data structures are very complex and there's a lot of intricate bits going on So I like to guarantee when function f is updating the database, no one else is and this is called serializable as a property and it's like one strongest thing you could ever say about the database and it costs some performance but it makes things so much easier to think about that I'm going to do the trade-off for a person.


[**slide**] 11:32
```
write transactions: db.Updates
```
        put := func(tx *bolt.Tx) error {
                bucket, err := tx.CreateBucket([]byte(“bukkit”))
                if err != nil {
                        return err
                }
                if err := bucket.Put([]byte(“answer”), []byte(“hello”)); err != nil {
                        return err
                }
                return nil
        }
        if err := db.Update(put); err != nil {
                log.Fatal(err)
        }
```
there is also db.Begin(writable bool), but you probably don’t want it
```


[**transcript**]
This is what a write transaction actually looks like, so you call the method db.Update and you give it a function that takes a uphold transaction and returned error and it may do things to the transaction to read and write the database, if it returns an error it gets rolled back, if it returns successful then the comments go through unless there's a disk error or something like that. So your store -- this is the actual error that comes whether it went through or not, but essentially if you say "oopsie I take that back" and returned an error, bolt undoes everything you did and to me-- I don't know about you, but I like writing these little function, I feel like I don't need a DSL embedded in string literals in my source code which is my annoyance with just about every possible thing is like I'm bringing another language inside this one language, why am I doing that? it just feels wrong, but this makes me feel I'm using go all the way.


[**slide**]  12:47
```
read transaction: db.View
```
        get := func(tx *bolt.Tx) error {
                bucket := tx.Bucket([]byte(“bukkit”))
                val := bucket.Get([]byte(“answer”))
                if val == nil {
                        // not found
                        return nil
        }
        if err  := db.View(get); err != nil {
                log.Fatal(err)
        }
```
don’t hold on to Get results
```


[**transcript**]
And read transaction look really similar but instead of update we say view, and the trick to know about this, is that all of the things you get out of bolt are only valid within the transaction, if you return from that function there, if it still holds the copy of the value you got from bucket.Get, things will go wrong, they will not be the same as when you last looked at it. So the lifetime of everything you receive from bolt is the life in the scope of the transaction and that's one of those things that I mean I made that mistake like, I discovered it in like two weeks ago? and I had made that mistake and it was on my code base for six months I think. so it's an easy mistake to make but this is a fundamental and anything that would do like, if you guys have used to syncing fork, for example, it's a way of life, Instead of churning in through garbage creating garbage for the garbage collector, it has a full of objects that are like managed a little bit more explicitly, it has the exact same thing, if you return an object into the pool it's essentially the same as the old C used at the free [inaudible 14:04], it's kinda the same thing except the language makes it a little bit safer but it's still like, that's -- that's one of the two difficult parts to make. so you kind of -- that's one sort of trap I'd say, that you kinda want to avoid.
[**slide**] 14:18
```
Bolt
```
key-value store
pure Go
library
single process
ordered 
b+tree
single writer
snapshot isolation...
```
[**transcript**]
So we talked about having a single writer, so what about the read side of the picture? so we can have at most one write transaction and at most one update code at a time. If another one tries-- it will wait until the first one complete, uh of course that makes no sense for read, you want to have as many reads going on as possible, that’s where all the performance comes from. even in sequel databases, you have this thing called isolation levels, which is like how the database vendor describes the behavior of the database for your queries. Snapshot isolation is pretty much the strongest possible thing you can say about a sequel database. So bolt is that level, so once again it’s really-- the most easy to reason about.


[**slide**] 15:11
```
snapshot isolation
demo: keep writing
```
        go func() {
                put := func(tx *bolt.Tx) error {
                        now := time.Now().Format(“15:04:05.000”)
                        fmt.Printf(“updating at\t%s\n”, now)
                        bucket := tx.Bucket([]byte(“bukkit”))
                        if err := bucket.Put([]byte(“clock”), []byte(now)); err != nil {
                                return err
                        }
                        return nil
                }
                for !done() {
                        if err := db.Update(put); err != nil {
                                log.Fatal(err)
                        }
                        time.Sleep(100 * time.Millisecond)
                }
        }()
```


[**transcript**]
Um, let’s not go with that yet..










[**slide**] 15:18
```
read transactions: db.View
```
        get := func(tx *bolt.Tx) error {
                bucket :=  tx.Bucket([]byte(“answer”))
                val := bucket. Get([]byte(“answer”))
                if val == nil {
                        // not found
                        return nil
                }
                fmt.Println(val)
                return nil
        }
        if err := db.View(get): err != nil {
                log.Fatal(err)
        }
```
don’t hold on to Get results
```
[**transcript**]
So what happens is when you start the transaction you call it db.View, you have a snapshot of the database as it was at the time your transaction started. And throughout this whole function, you would be reading that same snapshot. So if you read the same value again, it will still have the same old value.


[**slide**] 15:38
```
snapshot isolation
```
demo: keep writing
```
        go func() {
                put := func(tx *bolt.Tx) error {
                        now := time.Now().Format(“15:04:05.000”)
                        fmt.Printf(“updating at\t%s\”, now)
                        bucket := tx.Bucket([]byte(“bukkit”))
                        if err := bucket.Put([]byte(“clock”), []byte(now)); err != nil {
                                return err
                        }
                        return nil
                }
                for !done() {
                        if err := db.Update(put); err != nil {
                                log.Fatal(err)
                        }
                        time.Sleep(100 * time.Millisecond)
                }
        }()
```


[**transcript**]
So here’s an example of that, we have a writer. that will every 100 milliseconds write the current time of day into a single key called clock.


[**slide**]  15:50
```
snapshot isolation
```
demo: read after a delay
```
        time.Sleep(500 * time.Millisecond)
        var result []byte
        get := func(tx *bolt.Tx) error {
                fmt.Println(“reading…”)
                time.Sleep(500 *time.Millisecond)
                bucket := tx.Bucket([]byte(‘’bukkit”))
                clock := bucket.Get([]byte(“clock”))
                fmt.Printf(“observing %s\n”, clock)
                result = make([]byte, len(clock))
                copy(result, clock)
                return nil
        }
        if err := db.View(get); err != nil {
                log.Fatal(err)
        }
        fmt.Printf(“result %s\n”, result)
```
don’t hold on to Get results
don’t leak view transactions
```


[**transcript**]
And we delay a little and then start a transaction that also delays a little and then we say we are observing the code as of this time. So what happens when I run this, is that you get a bunch of updates-- so you have a bunch of updates going, 0 1 2 3 4, you have regular updates to the clock. And the transaction started here, by only ‘read the value’ here and it still seeing the old value. So it’s deceiving a snapshot of the database, which makes it easy to reason about things that don't mutate from underneath. The moment you started the transaction you have a consistent word.












[**slide**] 16:44
```
Bolt
```
key-value store
pure Go
library
single process
ordered 
b+tree
single writer
snapshot isolation
nested buckets...
```


[**transcript**]
And one of the things that make Bolt nicer than many other key-value stores is that you can use nested buckets to kind of structure your data. So instead of meaning to cram everything into a single main space where you know you need to encode into the key that-- this is my customer id-- and this is their whatever-- you can for example put customers, each customer into a separate bucket and underneath that-- customer for example--, every I don’t know -- maybe a product is a new bucket and under that is something, so you can structure the data with buckets that are essentially like folders where keys and value are like files and contents. Does that make sense?


[**slide**] 17:30
```
nested buckets: creating
```
        bucket, err := tx.CreateBucket([]byte(“bukkit”))
        if err !=nil  {
                return
        }
        sub, err := bucket.CreateBucket([]byte(“sub”))
        if err != nil {
        }
        deep, err := sub.CreateBucket([]byte(“deep”))
        if err != nil {
                return err
        }
```
scope for Bucket is the transaction; don’t store
```










[**transcript**]
So buckets are created with a code CreateBucket and you just get a reference to a bucket object and you can call methods on that to dig in deeper, whatever makes sense to you, for your application. And once again these things, the scope of the bucket value is the transaction, so you can’t hold onto it


[**slide**]  17:56
```
nested bucket: using
```
        bucket:= tx.Bucket9([]byte(“bukkit”))
        if bucket == nil {
                return errors.New(“bucket missing: bukkit”))
        }
        sub := bucket.Bucket([]byte(“sub”))
        if sub == nil {
                return errors.New(bucket missing: sub”))
        }
        deep := sub.Bucket([]byte(“deep”))
                return errors.New(bucket missing: deep”))
        }
```


[**transcript**]
And when you use them you just asked to get the bucket by this name and it’s outright nil or you get the actual bucket and you know if you call a bucket, the .Bucket method or the previous bucket, it’s a self bucket and so on, and there’s essentially no limit on how deep you can dig with this. There’s a little bit of overhead so you don’t easily want to create a bucket-- you know an object, a struct where every [inaudible] a key. that’s a little bit too much. So it’s-- once again Bolt is not a document store like, you know MongoDB or something like that.
[**slide**] 18:35


                                                   DB
    ↗   ⬆    ↖ ↖
                                        Tx        Tx    Tx  Tx
                                        ↑        ↑           
                                        Cursor  Bucket
                                                ↑
                                                 Bucket
                                                ↑
                                                 Cursor








[**transcript**]
So what you have is-- you have the database that has transactions pointing to it, and one of those transactions is the special write transaction, there can be only one write transaction at a time. And inside those transactions, you have cursors that are either doing-- its writing the top level or you have buckets, and bucket can point to other buckets and a cursor can be walking through the buckets. But this is essentially just a hierarchical tree accessing API.


[**slide**] 19:14 
```
Bolt
```
key-value store
pure Go
library
single process
ordered 
b+tree
single writer
snapshot isolation
nested buckets
zero-copy...
```


[**transcript**]
One really neat thing about Bolt is the way it’s implemented and this is logically the LMDB inspiration here.


[**slide**]  19:28
```
zero-copy
```
the database file is mmapped read-only
keys and values returned as [ ] byte point to this memory
only valid while the transaction is alive
read-only operations are just in-memory b+tree operations
kernel manages caching
```


[**transcript**]
So what they’ve done is that the whole database file is memory-mapped as a read-only area of memory, that file literally is in memory, as far as everything is concerned accessing that range of memory reads from the file without needing to do like syscall, without needing to do explicitly ask for file.Read or whatever. None of that sort of stuff-- that file is in memory and the read-only is what makes it safe, so you can’t accidentally write over it. So what Bolt does is let you do edit. When you commit a write transaction, when you call updates, it writes to the file but all the reading goes through the mmap. So if your machine has enough memory, it behaves as if it was an in-memory B+tree, which is kind of awesome when you think about it because usually in-memory data structures are fast and on-disk data structures are slow. But now if the active [inaudible 20:34] of your tree is small enough to fit into memory, it literally is an in-memory data structure. And there’s also no like double caching where the database library needs to maintain a row cache or something, and you know go through like, this is within the 1000 most requested rows let me hold it in my cache. All of that is done by the kernel and the big benefit of that is that first of all, the kernel, you can sort of assume that kernels are really well-written by smart people and second of all, the kernel has a lot of value, so it doesn’t matter what gets thrown out of memory, it could be possibly a database, it could be that JPEG that you haven’t served anyone for a while, you know file caches, it could be that popular executable that is only used once a month. Everything unnecessary gets thrown out of memory to make room for your database, or if something else more required your database gets thrown out of memory because it wasn’t the actively needed thing. So it’s just this weird thing where writing simple things creates something beautiful that performs well, and I really like that. That’s like when everything just clicks together. So this is why you can’t hold on to the get return value for example, because those bite-slices are literally pointing to this region of memory. And now remember it was a copy-on-write B tree with like, you know, space is being freed up and reduced. So there’s no guarantee when this bit of memory will be reused. Once you’ve ended the transaction, the transaction altered pinned to the old usage. But when the transaction in, that memory can be replaced with something else, some other content. So that to me like I knew LMDB before Bolt existed and I was trying to use LMDB for sequel writing and was like kinda into it, I like the idea and when I realized that Ben is writing Bolt to kind of imitate this one like, yup that’s a solid winner, I like that. It makes read operations really easy and lightweight. I think a read operation right now is something like full nanosecond, which is kinda hard to beat with anything that would like normally try to do things the other way


[**slide**] 23:10
```
read-only
```
don’t even try to mutate
```
        badIdea := func(tx *bolt.Tx) error {
                bucket := tx.Bucket([]byte(“bukkit”))
                val :=bucket.Get([]byte(“bork”))
                if val == nil {
                        return errors.New(“no answer”)
                }
                val[0] = 42
                return nil
        }
        if err :=db.Update(badIdea); err != nil {
                log.Fatal(err)
        }
```




[**transcript**]
So it’s what happens if you try to accidentally change the values returned to you, it crashes, crashes hard. So you get a paging value, this page is read-only that you tried to write into it. So, if you ever see this, you kind of know what you did. You probably failed to allocate a new buffer and copy this value into it.


[**slide**] 23:40
```
zero-copy downside
```
bolt database pages ae what kernel thinks they are
can’t use Snappy etc
kernel-level compression: chattr +c foo.bolt or dm-compress etc
works well with btrfs, often faster for all files
(and chattr +C for no CoW, it’s pointless and wasteful)
```


[**transcript**]
So, the downside of being memory-mapped is that the bolt database pages are exactly what the kernel thinks they are. I can say what's all of this, that is fungible, but Bolt database pages are what the kernel thinks they are. So there’s a no-load operation during which we could do decompression like LevelDB for example, when it writes up chunks of content to the disks it compresses it with snapping. And when it loads it from disk, it decompresses it with snapping. Now with Bolt, there’s no such thing as LWN, It’s just memory access. So there’s no convenient LWN to decompress anything. And I kind of think that that’s kind of unnecessary because these days a lot of stuff can compress files for you. So, for example with btrfs you can tell me this file should be stored compressed, and this actually a pretty good benchmark, if you say all of my files on this disk should all be compressed, it’s just gonna become faster. There’s also +C which says, “hey btrfs or any other copy-on-write part system, don’t give me the copy-on-write protections that you normally give files because this is an update in place file So, you know, avoid that overhead because it’s pointless, wasteful” and it’s the same thing that you do for example my virtual machine disk image that gets, you know, individual blocks written in it, It’s not like a single file that gets rewritten from scratch, so you can kind of optimize how things are done. If you’re not using one of the fanciest one-- one of the fancier file systems that have this stuff, You can also put the whole block device or something that compresses things on the fly. So dm compress is a block-level compression solution. But that kind of - is the downside is that you don’t have the flexibility of having custom code that launches things from this, because the loading has never occurred


[**slide**] 25:58
```
Bolt
```
key-value store
pure Go
library
single process
ordered 
b+tree
single writer
snapshot isolation
nested buckets
zero-copy
low-level...
```
[**transcript**]
Okay, Bolt is designed to do over, in case it hasn’t kinda come through yet, this is not really Bolt’s problem, this is the application's problem. A lot of things, these things actually kinda do-- kind of-- even the LWN in the application, because then you can do things in a way that is small and takes advantage of what you know in the application.


[**slide**] 26:16
```
low-level
```
no schema, indices, or queries
sequential transaction batching belongs in the application
(because off error handling)
(but see on-going DB.Batch work for batching concurrent transactions)


no Write-ahead Log (WAL) here
```


[**transcript**]
So Bolt doesn’t do things like schemas, indices, queries there is no DSL, there’s no like filtering off, giving me all the keys that have the value “hello”, there is no such thing. You are-you walk through all of them, or you write an index that helps you do that thing. And what I kinda wanna see is a little ecosystem of indexing built on top of the interface. So, for example, there’s no reason why there’s that project “bleve”? I’m not sure how you pronounce that, that’s essentially a search engine that the guts of a search engine written in go. There’s no reason why you couldn’t take one of those things that are the very generic engine and even a stolen index in Bolt or whatever and updated it from your application, but that’s kind of outside of the problem domain of Bolt itself, and that’s actually really relieving because you know that Bolt is pretty much gonna be what it is, like it’s not something that is gonna alive in the next two years and be something else. I kind of like this, very stable WAL. We’re still talking about some additions, some minor API changes, but we feel like it’s pretty much like definitely a 1.0 gap. So a good example of what is not a Bolt problem is a sequential transaction batch. So the idea of in the same disc sync doing multiple updates in that one synchronization operation, and this is something that a sequel database would do for you. It will combine multiple queries that you send in-- not queries but multiple transactions that you send in. It will combine those and do multiple at once. And this sort of stuff kinda isn’t Bolt’s problem, especially when it’s sequential, so like imagine the idea of you wanting to import a file that has 10 million lines and you create a key for each line. Now the trivial way of doing it is to read a line, pass it, update, set one key, read the next line, pass it, update, set key. The problem with that is that you’re doing two full disc syncs for every line, that’s 20 million disc sync which is gonna take 40 hours, that’s not good. So what you do instead is you want a batch Bolt an update at once. So let’s say you pass 1000 lines, do a transaction that writes those out and you can make that be streaming so you’re not holding them on memory on your application-level and the reason why that belongs in the application is that you don’t whether it quit until the end of those 1000 lines. So we can have the same simple API when you say put this, tell me if it works if you’re trying to sequentially batch them, instead we need to actually do this batch and then check whether it worked. We are doing a thing that we called DB.batch which is taking concurrent operations and merging those into a single transaction. So for example, if you have a web app, an API, or some sort that takes incoming requests and does something to the database as a result of those requests, it with this thing,  let’s say anything going in within the same 10 milliseconds up to a limit of 1000 operations gets put in one transaction. So we have that, but it is now a concurrent transaction. not like; let’s import that one file out of the operation. 


[**slide**] 29:57 


```
Bolt
```
key-value store
pure Go
library
single process
ordered 
b+tree
single writer
snapshot isolation
nested buckets
zero-copy
low-level
```


[**transcript**]
So, Bolt is a key-value store, so no relational databases, no sequel, no queries, nothing like that. Simple, it is speaking of Pure Go, it’s a pretty project, it is a simple project that is surprisingly a really small project. It’s a library, it’s not a service, it’s not a demand you run, it’s something that becomes part of your application. It’s meant to be-- the Bolt databases meant to be used from a single process. So, you really own the file that you have all, nobody is using it at the same time.  So, for example, a sequel line will let you processes, open the same file, and it [inaudible]. Bolt kind of try [inaudible]. The key-value store is ordered and you can use that to do interesting things, like for example, I stored directory listing in Bolt, and I know I can work them in alphabetical order just because it’s an ordered key-value store. It is internally a B+ tree which kind of works, it’s like one of the still of the art data structures. I mean you can do better but not much better. So B+tree is pretty neat. It is based on the idea of that being at most the single writer at a time to the database, and that makes application one-degree simpler to write. By reading inside you get snapshot isolations which means whenever you’re reading from the database, you have a consistent WAL view. You don’t have to deal with things changing from underneath you and you can use nested bucket structure in your application, Bolt is internally built using mmap which is a really neat trick, and works really well for Bolt. And it’s kind of a low-level tool. So it may not do what you want directly but it is really like workable when you think of it as a library, You need to kinda embrace the fact that it is just one of the things that your application uses and I really enjoy using it and it has been a great match for my application. I would like to hear-- maybe later once all the talks are done, come to me to talk about, like, I have this and this thing, could Bolt be used for this and this, That’s all the conversations are-- in my mind really useful because you get to kinda evaluate, does this fit this at all? does my application benefit from it at all? So are there any questions? Yap, in the back. Is there a mechanism to do back up and restore? Good question. Yes, of course, there is. So when-- if you read the file, the database file from outside the process, it can change from underneath you. So, that’s not a safe way to do the backup. So what Bolt actually has-- in the library is you can call a function and it will essentially write a copy of the database that is a snapshot to your thingie. I forget whether it is-- I think it handed an io.Writer and it just writes to copy of the database that i.o written. Also if you have any solutions at all that can give you an atomic snapshot which for example [inaudible]. That’s the perfectly valid way of taking a backup because Bolt is a crash-safe. So if-- even if you kind of catch it in the middle of doing changes, you rolled back to the previous step. Okay, Yea behind. Under what circumstances do you reestablish the memory map? Reestablish the mmap-- for example behind the line file [inaudible] how do you get memory-mapped access--? Yes, so mmap is a fixed range of memory. You can’t, like, change the size of it without mmapping something else. So when the database runs out of free pages and needs another free page, what it does is it increments the size of the file and re-mmaps it with the bigger size, and now it has a bunch of free pages. Does that make sense? Yea. Am I (..34:14). What basically you’re holding onto in your pointers? That was the whole point of like-- the values that get returned to you are only valid within the transaction, What about the performance? Since it’s just a single writer, can it become a bottleneck in my applications? So um, that thing with a single write, if you remember I said about the copy-on-write B tree-- there we go. So, to flush out a comment, Bolt needs to write all the joint pages, make sure they are stable on the disk, then write a new metadata page, and make sure that is stable on the disk. And then it says we’re done. So, every transaction, every update call costs 2 disk syncs. And that kind of [inaudible]. Honestly, how much data you write is almost irrelevant, because the syncs themselves are the cost. Now an SSD will do like some hundreds of syncs per second. So that’s your limit, and that’s what the motivation of the batching mechanism comes in. The read side, when things go well, it’s literally an in-memory B+ tree. And I am getting timing of nanoseconds. So, Bolt is definitely meant to be a read-heavy solution. And the way I think of it is-- if you’re building something that is not-- let me think about that. So, Bolt is not a network service in itself. So, there’s no such thing as multiple things coming to Bolt directly. So there’s always more application mediating that conversation from any clients you have. So, in that sort of situation your application is the perfect place to do some sort of batching and gets-- it puts doing [inaudible] of the error handling what if the batch fails, you know that sort of stuff. A DB batch for your average network application that does something in the style of [inaudible] coming in or requests coming in. DB batch is pretty much what you need. It’s close enough to the perfect solution. In my mind, it is completely an acceptable solution.
