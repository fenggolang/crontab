golang的chan的设计很好：无缓存chan相当于读写锁，有缓存chan相当于有锁队列。
