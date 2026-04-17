import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { GoApiError, goFetch } from './go-client'

describe('goFetch', () => {
  const originalEnv = process.env.BACKEND_GO_URL

  beforeEach(() => {
    process.env.BACKEND_GO_URL = 'https://api5g.numbro.app'
    vi.stubGlobal('fetch', vi.fn())
  })

  afterEach(() => {
    process.env.BACKEND_GO_URL = originalEnv
    vi.unstubAllGlobals()
  })

  it('returns parsed JSON on successful GET', async () => {
    const mockData = { id: '123', status: 'ok' }
    vi.mocked(fetch).mockResolvedValueOnce(
      new Response(JSON.stringify(mockData), { status: 200 })
    )

    const result = await goFetch<typeof mockData>('/v1/consumer-units')
    expect(result).toEqual(mockData)
  })

  it('throws GoApiError with correct status on non-2xx response', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(
      new Response(JSON.stringify({ error: 'not found' }), { status: 404 })
    )
    await expect(goFetch('/v1/consumer-units/999')).rejects.toThrow(GoApiError)

    vi.mocked(fetch).mockResolvedValueOnce(
      new Response(JSON.stringify({ error: 'not found' }), { status: 404 })
    )
    await expect(goFetch('/v1/consumer-units/999')).rejects.toMatchObject({
      status: 404,
      path: '/v1/consumer-units/999',
    })
  })

  it('throws GoApiError with status 504 on network timeout', async () => {
    vi.mocked(fetch).mockRejectedValueOnce(
      Object.assign(new DOMException('signal timed out', 'AbortError'))
    )

    await expect(goFetch('/v1/consumer-units', { timeoutMs: 100 })).rejects.toMatchObject({
      status: 504,
    })
  })

  it('throws configuration error when BACKEND_GO_URL is missing', async () => {
    delete process.env.BACKEND_GO_URL

    await expect(goFetch('/v1/consumer-units')).rejects.toThrow('BACKEND_GO_URL')
  })
})
