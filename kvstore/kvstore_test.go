package kvstore

 

import (

    "fmt"

    "sync"

    "testing"

)

 

// TestBasicOperations tests Set, Get, and Delete in a single-threaded environment.

func TestBasicOperations(t *testing.T) {

    kv := NewKVStore()

 

    // Test Get on an empty store

    _, found := kv.Get("name")

    if found {

        t.Error("Expected 'name' to not be found in empty store")

    }

 

    // Test Set and Get

    kv.Set("name", "Ben")

    val, found := kv.Get("name")

    if !found {

        t.Error("Expected to find 'name'")

    }

    if val != "Ben" {

        t.Errorf("Expected value 'Ben', got %q", val)

    }

 

    kv.Set("name", "Bob")

    val, found = kv.Get("name")

    if !found {

        t.Error("Expected to find 'name' after overwrite")

    }

    if val != "Bob" {

        t.Errorf("Expected overwritten value 'Bob', got %q", val)

    }

 

    // Test Delete

    kv.Delete("name")

    _, found = kv.Get("name")

    if found {

        t.Error("Expected 'name' to be deleted")

    }

}

 

func TestConcurrentAccess(t *testing.T) {

    kv := NewKVStore()

    var wg sync.WaitGroup

 

    numWorkers := 50

    iterations := 100

 

    // Launch concurrent writers

    for i := 0; i < numWorkers; i++ {

        wg.Add(1)

        go func(workerID int) {

            defer wg.Done()

            for j := 0; j < iterations; j++ {

                key := fmt.Sprintf("key-%d-%d", workerID, j)

                kv.Set(key, "value")

            }

        }(i)

    }

 

    // Launch concurrent readers

    for i := 0; i < numWorkers; i++ {

        wg.Add(1)

        go func(workerID int) {

            defer wg.Done()

            for j := 0; j < iterations; j++ {

                key := fmt.Sprintf("key-%d-%d", workerID, j)

                // We don't care if it's found or not (since writers are racing with readers),

                // we just want to ensure reading doesn't panic or cause a data race.

                _, _ = kv.Get(key)

            }

        }(i)

    }

 

    // Wait for all goroutines to finish

    wg.Wait()

 

    value, found := kv.Get("key-0-0")

    if !found {

        t.Error("expected key-0-0 to exist")

    }

    if value != "value" {

        t.Errorf("expected value, got %q", value)

    }

}