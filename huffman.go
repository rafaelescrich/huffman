package huffman

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// HNode is a node in a Huffman Tree
type HNode struct {
	value     byte
	frequency int
	left      *HNode
	right     *HNode
}

// HTree holds a head node of Huffman Tree
type HTree struct {
	reader *bufio.Scanner
	root   *HNode
}

// New initializes a new instance
func New(reader io.Reader) *HTree {
	ht := &HTree{bufio.NewScanner(reader), nil}
	if err := ht.generateHuffmanTree(); err != nil {
		return nil
	}
	return ht
}

// generateTree generates a huffman tree
func (h *HTree) generateHuffmanTree() error {
	hash := make(map[byte]int)
	for h.reader.Scan() {
		line := h.reader.Text()
		fmt.Println(line)
		for _, l := range line {
			hash[byte(l)] += 1
		}
	}

	if len(hash) < 1 {
		return errors.New("huffman: no huffman nodes generated")
	}

	var nodes []*HNode
	for v, f := range hash {
		nodes = append(nodes, &HNode{v, f, nil, nil})
	}
	fmt.Println("Hash: ", hash)

	h.sort(nodes, 0, len(nodes)-1)

	for len(nodes) > 1 {
		fmt.Println("Nodes: ", nodes)
		newNode := &HNode{byte('*'), nodes[0].frequency + nodes[1].frequency, nodes[0], nodes[1]}
		nodes[1] = newNode
		nodes = nodes[1:]
		h.sort(nodes, 0, len(nodes)-1)
		fmt.Println("Nodes After: ", nodes)
	}

	h.root = nodes[0]
	nodes[0] = nil
	fmt.Println("Root Node: ", h.root.frequency, h.root.value)
	return nil
}

// Encode encodes the given string
func (h *HTree) Encode() error {
	if h.root == nil {
		return errors.New("huffman: can't encode, no huffman tree found.")
	}

	reader := bufio.NewScanner(strings.NewReader("Eerie eyes seen near lake."))
	var buffer bytes.Buffer
	for reader.Scan() {
		line := reader.Text()
		fmt.Println(line)
		for _, l := range line {
			var data []byte
			data = h.root.encode(byte(l), data)
			buffer.Write(data)
		}
	}

	fp, err := os.Create("encoded_file.txt")
	if err != nil {
		return err
	}

	if _, err = fp.Write(buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

func (hn *HNode) encode(marker byte, data []byte) []byte {
	if hn != nil {
		if hn.value == marker {
			return data
		}

		if hn.left != nil {
			data = append(data, byte('0'))
			hn.left.encode(marker, data)
		}

		if hn.right != nil {
			data = append(data, byte('1'))
			hn.right.encode(marker, data)
		}
	}
	return data
}

// Decode decodes the encoded data
func (h *HTree) Decode(encoded string) []byte {
	if h.root == nil {
		return nil
	}

	var decoded []byte
	node := h.root
	for _, e := range encoded {
		if node.left == nil && node.right == nil {
			decoded = append(decoded, node.value)
			node = h.root
		}
		if byte(e) == byte('0') {
			node = node.left
		}
		if byte(e) == byte('1') {
			node = node.right
		}
	}
	return decoded
}

// Traverse helps walkthrough huffman tree
func (h *HTree) Traverse() {
	if h.root != nil {
		h.root.traverse()
	}
}

// sort implementating Quicksort for sorting list of huffman nodes
func (h *HTree) sort(nodes []*HNode, start, end int) {
	if start < end {
		p := h.partition(nodes, start, end)
		h.sort(nodes, 0, p-1)
		h.sort(nodes, p, end)
	}
}

// part of Quicksort
func (h *HTree) partition(nodes []*HNode, start, end int) int {
	pivot := nodes[(start+end)/2].frequency
	i, j := start, end

	for i <= j {
		for nodes[i].frequency < pivot {
			i++
		}

		for nodes[j].frequency > pivot {
			j--
		}

		if i <= j {
			nodes[i], nodes[j] = nodes[j], nodes[i]
			i++
			j--
		}
	}

	return i
}

func (hn *HNode) traverse() {
	if hn != nil {
		hn.left.traverse()
		fmt.Println(hn.value, hn.frequency)
		hn.right.traverse()
	}
}
