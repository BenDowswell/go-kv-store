// Store is the interface implemented by all KV store backends.

package kvstore

 

type Store interface {

    Get(key string) (string, bool)

    Set(key, value string)

    Delete(key string)

}