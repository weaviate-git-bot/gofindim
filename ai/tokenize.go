package ai

import (
	// Import other libraries

	"fmt"
	"os"

	"github.com/daulet/tokenizers"
)

func Tokenize(text string) ([]uint32, error) {

	file := "ai_assets/tokenizer.json"
	_, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	// Check path
	tk, err := tokenizers.FromFile(file)
	if err != nil {
		return nil, err
	}

	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	} else {
		fmt.Println("File exists")
	}

	tokens := tk.Encode(text, false)

	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// I do have all the files from hugging face like tokenizer_config, tokenizer, and vocab. How do I use those to tokenize my text? I am using `"github.com/tiktoken-go/tokenizer`. Is this compatible with CLIP? I'm using "
// CLIP-ViT-L-14-DataComp.XL-s13B-b90K" Just answer the question without code.
