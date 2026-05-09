import dayjs from 'dayjs'

export function formatDateTime(value) {
  if (!value) {
    return '—'
  }

  return dayjs(value).isValid() ? dayjs(value).format('DD.MM.YYYY HH:mm:ss') : value
}

export function formatJsonPreview(value) {
  if (value === null || value === undefined) {
    return '—'
  }

  return JSON.stringify(value)
}

export function resolveRelationLabel({
  value,
  source,
  optionValueKey,
  getOptionLabel,
  optionsMap,
}) {
  if (value === null || value === undefined || value === 0) {
    return '—'
  }

  const items = optionsMap[source] ?? []
  const matched = items.find((item) => item[optionValueKey] === value)

  if (!matched) {
    return `#${value}`
  }

  return getOptionLabel ? getOptionLabel(matched) : matched[optionValueKey]
}
