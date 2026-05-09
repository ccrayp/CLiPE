import { http } from './http'

function unwrapData(response) {
  return response?.data?.data ?? {}
}

function compactPayload(payload) {
  const result = {}

  Object.entries(payload ?? {}).forEach(([key, value]) => {
    if (value === '' || value === undefined) {
      return
    }

    result[key] = value
  })

  return result
}

export function createCrudApi({
  basePath,
  listKey,
  buildUpdatePath = (id) => `${basePath}/${id}`,
  buildDeletePath = (id) => `${basePath}/${id}`,
}) {
  return {
    async search(filters = {}, limit = 10, offset = 0) {
      const response = await http.post(
        `${basePath}/search?limit=${limit}&offset=${offset}`,
        compactPayload(filters),
      )

      const data = unwrapData(response)

      return {
        items: data[listKey] ?? [],
        count: Number(data.count ?? 0),
        limit: Number(data.limit ?? limit),
        offset: Number(data.offset ?? offset),
      }
    },

    async create(payload) {
      const response = await http.post(basePath, compactPayload(payload))
      return unwrapData(response)
    },

    async update(id, payload) {
      const response = await http.put(
        buildUpdatePath(id),
        compactPayload(payload),
      )
      return unwrapData(response)
    },

    async remove(id) {
      const response = await http.delete(buildDeletePath(id))
      return unwrapData(response)
    },
  }
}
