## Design Decisions
- Memcached is used to cache ID along with timestamp everytime there is a new request to `/api/verve/accept` end point. This choice is made since at any point, we need a maximum of last 2 minutes keys to be stored with timestamp in cache.
Note: This was not the best choice due to the Memcached limitation of not being able to read all keys in cache at once. I was late to realise this and handled it differently using additional `keyList` in the cache. Alternatively we can use Redis which supports querying cache for data.

- In `/api/verve/accept` endpoint, used go routines to parallelly write to memcached and make endpoint api call after ID validation to reduce the overall API processing time.

- Used singleton pattern to initialize memcached client and reuse same client all the times to save on initialization times.

- Scheduled a periodic job that counts and logs the last minute unique id counts to a log file and clears the last minute cache afterwards.

- I have also implemented `/endpoint` APIs for using as target endpoint in `/api/verve/accept` call. 