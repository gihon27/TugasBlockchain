// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Paket konsensus mengimplementasikan mesin konsensus Ethereum yang berbeda.
package consensus

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

// ChainHeaderReader mendefinisikan kumpulan kecil metode yang diperlukan untuk mengakses lokal
// blockchain selama verifikasi header.
type ChainHeaderReader interface {
	// Config retrieves the blockchain's chain configuration.
	Config() *params.ChainConfig

	// CurrentHeader retrieves the current header from the local chain.
	CurrentHeader() *types.Header

	// GetHeader retrieves a block header from the database by hash and number.
	GetHeader(hash common.Hash, number uint64) *types.Header

	// GetHeaderByNumber retrieves a block header from the database by number.
	GetHeaderByNumber(number uint64) *types.Header

	// GetHeaderByHash retrieves a block header from the database by its hash.
	GetHeaderByHash(hash common.Hash) *types.Header

	// GetTd retrieves the total difficulty from the database by hash and number.
	GetTd(hash common.Hash, number uint64) *big.Int
}

// ChainReader mendefinisikan kumpulan kecil metode yang diperlukan untuk mengakses lokal
// blockchain selama verifikasi header dan/atau uncle verification.
type ChainReader interface {
	ChainHeaderReader

	// GetBlock retrieves a block from the database by hash and number.
	GetBlock(hash common.Hash, number uint64) *types.Block
}

// Engine adalah mesin konsensus agnostik algoritma.
type Engine interface {
	// Penulis mengambil alamat Ethereum dari akun yang mencetak blok yang diberikan, 
	// yang mungkin berbeda dari basis koin header jika mesin konsensus didasarkan pada tanda tangan.
	Author(header *types.Header) (common.Address, error)

	// VerifyHeader memeriksa apakah header sesuai dengan aturan konsensus dari mesin yang diberikan. 
	// Memverifikasi segel dapat dilakukan secara opsional di sini, atau secara eksplisit melalui metode VerifySeal.
	VerifyHeader(chain ChainHeaderReader, header *types.Header, seal bool) error

	// VerifyHeaders mirip dengan VerifyHeader, tetapi memverifikasi sekumpulan header secara bersamaan. 
	// Metode ini mengembalikan saluran keluar untuk membatalkan operasi dan saluran hasil untuk mengambil 
	// verifikasi asinkron (urutan adalah urutan irisan input).
	VerifyHeaders(chain ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error)

	// VerifyUncles memverifikasi bahwa paman blok yang diberikan 
	// sesuai dengan aturan konsensus dari mesin yang diberikan.
	VerifyUncles(chain ChainReader, block *types.Block) error

	// Siapkan menginisialisasi bidang konsensus dari header blok sesuai dengan aturan mesin tertentu. 
	// Perubahan dijalankan sebaris.
	Prepare(chain ChainHeaderReader, header *types.Header) error

	// Finalize menjalankan modifikasi status pasca-transaksi apa pun (misalnya hadiah blok) tetapi tidak merakit blok.
	// Catatan: Header blok dan database negara bagian mungkin diperbarui untuk mencerminkan aturan konsensus 
	// apa pun yang terjadi pada finalisasi (misalnya, hadiah blok).
	Finalize(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header)

	// FinalizeAndAssemble menjalankan modifikasi status pasca-transaksi (misalnya hadiah blok) dan merakit blok terakhir.
	// Catatan: Header blok dan database negara bagian mungkin diperbarui untuk 
	// mencerminkan aturan konsensus apa pun yang terjadi pada finalisasi (misalnya, hadiah blok).
	FinalizeAndAssemble(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error)

	// Seal menghasilkan permintaan penyegelan baru untuk blok input yang diberikan dan mendorong hasilnya ke saluran yang diberikan.
	// Catatan, metode ini segera kembali dan akan mengirimkan hasil async. 
	// Lebih dari satu hasil juga dapat dikembalikan tergantung pada algoritma konsensus.
	Seal(chain ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error

	// SealHash mengembalikan hash dari sebuah blok sebelum disegel.
	SealHash(header *types.Header) common.Hash

	// CalcDifficulty adalah algoritma penyesuaian kesulitan. Ini mengembalikan kesulitan
	// yang harus dimiliki oleh blok baru.
	CalcDifficulty(chain ChainHeaderReader, time uint64, parent *types.Header) *big.Int

	// API mengembalikan API RPC yang disediakan mesin konsensus ini.
	APIs(chain ChainHeaderReader) []rpc.API

	// Tutup mengakhiri semua utas latar belakang yang dikelola oleh mesin konsensus.
	Close() error
}

// PoW adalah mesin konsensus berdasarkan proof-of-work.
type PoW interface {
	Engine

	// Hashrate mengembalikan hashrate penambangan saat ini dari mesin konsensus PoW.
	Hashrate() float64
}
