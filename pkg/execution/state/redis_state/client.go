package redis_state

import (
	"context"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/redis/rueidis"
)

const StateDefaultKey = "estate"
const QueueDefaultKey = "queue"

type FunctionRunStateClient struct {
	kg            RunStateKeyGenerator
	client        RetriableClient
	unshardedConn RetriableClient
	isSharded     IsShardedFn
}

func (f *FunctionRunStateClient) KeyGenerator() RunStateKeyGenerator {
	return f.kg
}

func (f *FunctionRunStateClient) Client(ctx context.Context, accountId uuid.UUID, runId ulid.ULID) (RetriableClient, bool) {
	if f.isSharded(ctx, accountId, runId) {
		return f.client, true
	}
	return f.unshardedConn, false
}

func (f *FunctionRunStateClient) ForceShardedClient() RetriableClient {
	return f.client
}

func NewFunctionRunStateClient(r rueidis.Client, u *UnshardedClient, stateDefaultKey string, isSharded IsShardedFn) *FunctionRunStateClient {
	return &FunctionRunStateClient{
		kg:            &runStateKeyGenerator{stateDefaultKey: stateDefaultKey},
		client:        newRetryClusterDownClient(r),
		unshardedConn: NewNoopRetriableClient(u.unshardedConn),
		isSharded:     isSharded,
	}
}

type ShardedClient struct {
	fnRunState *FunctionRunStateClient
}

type IsShardedFn func(ctx context.Context, accountId uuid.UUID, runId ulid.ULID) bool

func AlwaysShard(ctx context.Context, accountId uuid.UUID, runId ulid.ULID) bool {
	return true
}

func NeverShard(ctx context.Context, accountId uuid.UUID, runId ulid.ULID) bool {
	return false
}

type ShardedClientOpts struct {
	UnshardedClient        *UnshardedClient
	FunctionRunStateClient rueidis.Client
	StateDefaultKey        string
	FnRunIsSharded         IsShardedFn
}

func NewShardedClient(opts ShardedClientOpts) *ShardedClient {
	return &ShardedClient{
		fnRunState: NewFunctionRunStateClient(opts.FunctionRunStateClient, opts.UnshardedClient, opts.StateDefaultKey, opts.FnRunIsSharded),
	}
}

func (s *ShardedClient) FunctionRunState() *FunctionRunStateClient {
	return s.fnRunState
}

type PauseClient struct {
	kg          PauseKeyGenerator
	unshardedRc rueidis.Client
}

func (p *PauseClient) KeyGenerator() PauseKeyGenerator {
	return p.kg
}

func (p *PauseClient) Client() rueidis.Client {
	return p.unshardedRc
}

func NewPauseClient(r rueidis.Client, stateDefaultKey string) *PauseClient {
	return &PauseClient{
		kg:          pauseKeyGenerator{stateDefaultKey: stateDefaultKey},
		unshardedRc: r,
	}
}

type QueueClient struct {
	kg          QueueKeyGenerator
	unshardedRc rueidis.Client
}

func (q *QueueClient) KeyGenerator() QueueKeyGenerator {
	return q.kg
}

func (q *QueueClient) Client() rueidis.Client {
	return q.unshardedRc
}

func NewQueueClient(r rueidis.Client, queueDefaultKey string) *QueueClient {
	return &QueueClient{
		kg:          queueKeyGenerator{queueDefaultKey: queueDefaultKey, queueItemKeyGenerator: queueItemKeyGenerator{queueDefaultKey: queueDefaultKey}},
		unshardedRc: r,
	}
}

type BatchClient struct {
	kg          BatchKeyGenerator
	unshardedRc rueidis.Client
}

func (b *BatchClient) KeyGenerator() BatchKeyGenerator {
	return b.kg
}

func (b *BatchClient) Client() rueidis.Client {
	return b.unshardedRc
}

func NewBatchClient(r rueidis.Client, queueDefaultKey string) *BatchClient {
	return &BatchClient{
		kg:          batchKeyGenerator{queueDefaultKey: queueDefaultKey, queueItemKeyGenerator: queueItemKeyGenerator{queueDefaultKey: queueDefaultKey}},
		unshardedRc: r,
	}
}

type DebounceClient struct {
	kg          DebounceKeyGenerator
	unshardedRc rueidis.Client
}

func (d *DebounceClient) KeyGenerator() DebounceKeyGenerator {
	return d.kg
}

func (d *DebounceClient) Client() rueidis.Client {
	return d.unshardedRc
}

func NewDebounceClient(r rueidis.Client, queueDefaultKey string) *DebounceClient {
	return &DebounceClient{
		kg:          debounceKeyGenerator{queueDefaultKey: queueDefaultKey, queueItemKeyGenerator: queueItemKeyGenerator{queueDefaultKey: queueDefaultKey}},
		unshardedRc: r,
	}
}

type GlobalClient struct {
	kg          GlobalKeyGenerator
	unshardedRc rueidis.Client
}

func (g *GlobalClient) KeyGenerator() GlobalKeyGenerator {
	return g.kg
}

func (g *GlobalClient) Client() rueidis.Client {
	return g.unshardedRc
}

func NewGlobalClient(r rueidis.Client, stateDefaultKey string) *GlobalClient {
	return &GlobalClient{
		kg:          globalKeyGenerator{stateDefaultKey: stateDefaultKey},
		unshardedRc: r,
	}
}

type UnshardedClient struct {
	unshardedConn rueidis.Client

	pauses   *PauseClient
	queue    *QueueClient
	batch    *BatchClient
	debounce *DebounceClient
	global   *GlobalClient
}

func (u *UnshardedClient) Pauses() *PauseClient {
	return u.pauses
}

func (u *UnshardedClient) Queue() *QueueClient {
	return u.queue
}

func (u *UnshardedClient) Batch() *BatchClient {
	return u.batch
}

func (u *UnshardedClient) Debounce() *DebounceClient {
	return u.debounce
}

func (u *UnshardedClient) Global() *GlobalClient {
	return u.global
}

func NewUnshardedClient(r rueidis.Client, stateDefaultKey, queueDefaultKey string) *UnshardedClient {
	return &UnshardedClient{
		pauses:        NewPauseClient(r, stateDefaultKey),
		queue:         NewQueueClient(r, queueDefaultKey),
		batch:         NewBatchClient(r, queueDefaultKey),
		debounce:      NewDebounceClient(r, queueDefaultKey),
		global:        NewGlobalClient(r, stateDefaultKey),
		unshardedConn: r,
	}
}