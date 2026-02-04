package main

type TransactionInput struct {
	PrevTxID     []byte
	PrevOutIndex int
	Signature    []byte
	PubKey       []byte
}
