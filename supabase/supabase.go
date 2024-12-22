package supabase1

import (
	"fmt"
	"github.com/supabase-community/auth-go"
	

	
)

func NewSupabaseClient(url string, key string) auth.Client  {
    if url == "" || key == "" {
       fmt.Println("url and key cannot be empty")
	   return nil;
    }
	client := auth.New(url,key);

    return client;
}

