/* eslint-disable react-hooks/set-state-in-effect */
import { useEffect, useMemo, useRef, useState } from 'react'
import { Alert, Button } from 'react-bootstrap'
import { useSearchParams } from 'react-router-dom'
import { apiMap } from '../api'
import { getApiErrorMessage } from '../api/http'
import { ConfirmModal } from '../components/common/ConfirmModal'
import { DataTable } from '../components/common/DataTable'
import { EntityFormModal } from '../components/common/EntityFormModal'
import { JsonViewerModal } from '../components/common/JsonViewerModal'
import { PageHeader } from '../components/common/PageHeader'
import { SearchPanel } from '../components/common/SearchPanel'
import { IoMdAdd } from "react-icons/io";

const DEFAULT_PAGE_SIZE = 10
const SEARCH_DEBOUNCE_MS = 350

function buildInitialFilters(fields, searchParams) {
  return fields.reduce((accumulator, field) => {
    accumulator[field.name] = searchParams.get(field.name) ?? ''
    return accumulator
  }, {})
}

function normalizeFilters(fields, values) {
  return fields.reduce((accumulator, field) => {
    const rawValue = values[field.name]

    if (rawValue === '') {
      return accumulator
    }

    if (field.type === 'number' || field.type === 'select') {
      accumulator[field.name] = Number(rawValue)
      return accumulator
    }

    if (field.type === 'triBoolean') {
      accumulator[field.name] = rawValue === 'true'
      return accumulator
    }

    accumulator[field.name] = rawValue
    return accumulator
  }, {})
}

function collectSources(config) {
  const allFields = [...config.filters, ...config.formFields, ...config.columns, ...(config.detailFields ?? [])]

  return [...new Set(allFields.filter((field) => field.source).map((field) => field.source))]
}

export function EntityPage({ config }) {
  const [searchParams, setSearchParams] = useSearchParams()
  const [filters, setFilters] = useState(() => buildInitialFilters(config.filters, searchParams))
  const [debouncedFilters, setDebouncedFilters] = useState(filters)
  const [rows, setRows] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(() => {
    const pageParam = searchParams.get('page')
    return pageParam ? Number(pageParam) : 1
  })
  const [pageSize, setPageSize] = useState(() => {
    const limit = Number(searchParams.get('limit') ?? DEFAULT_PAGE_SIZE)
    return [10, 20, 50, 100].includes(limit) ? limit : DEFAULT_PAGE_SIZE
  })
  const [loading, setLoading] = useState(true)
  const [isFetching, setIsFetching] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [optionsMap, setOptionsMap] = useState({})
  const [modalState, setModalState] = useState({
    show: false,
    mode: 'create',
    record: null,
  })
  const [deleteRecord, setDeleteRecord] = useState(null)
  const [jsonState, setJsonState] = useState({
    show: false,
    title: '',
    value: null,
  })
  const rowsCountRef = useRef(0)
  const shouldResetPageRef = useRef(false) // Флаг для сброса страницы

  const initialFilters = useMemo(
    () => buildInitialFilters(config.filters, searchParams),
    [config.filters, searchParams],
  )

  // Сброс страницы только когда меняются фильтры (не при пагинации)
  useEffect(() => {
    // Синхронизируем фильтры из URL
    setFilters(initialFilters)
    setDebouncedFilters(initialFilters)

    // Синхронизируем страницу из URL
    const pageParam = searchParams.get('page')
    const nextPage = pageParam ? Number(pageParam) : 1
    setPage(nextPage)

    // Синхронизируем размер страницы из URL
    const limitParam = Number(searchParams.get('limit') ?? DEFAULT_PAGE_SIZE)
    const nextPageSize = [10, 20, 50, 100].includes(limitParam)
      ? limitParam
      : DEFAULT_PAGE_SIZE

    setPageSize(nextPageSize)
  }, [initialFilters, config.key, searchParams])

  // Дебаунс фильтров - НЕ сбрасываем страницу здесь!
  useEffect(() => {
    const timer = window.setTimeout(() => {
      setDebouncedFilters(filters)
    }, SEARCH_DEBOUNCE_MS)

    return () => window.clearTimeout(timer)
  }, [filters])

 

  // Синхронизация URL с состоянием
  useEffect(() => {
    const nextParams = new URLSearchParams()

    Object.entries(filters).forEach(([key, value]) => {
      if (value !== '') {
        nextParams.set(key, value)
      }
    })

    nextParams.set('limit', String(pageSize))
    if (page > 1) {
      nextParams.set('page', String(page))
    }

    setSearchParams(nextParams, { replace: true })
  }, [filters, pageSize, page, setSearchParams])

  useEffect(() => {
    let mounted = true

    async function loadOptions() {
      const sources = collectSources(config)

      if (sources.length === 0) {
        return
      }

      try {
        const entries = await Promise.all(
          sources.map(async (source) => {
            const result = await apiMap[source].search({}, 500, 0)
            return [source, result.items]
          }),
        )

        if (mounted) {
          setOptionsMap(Object.fromEntries(entries))
        }
      } catch (requestError) {
        if (mounted) {
          setError(getApiErrorMessage(requestError, 'Не удалось загрузить связанные данные'))
        }
      }
    }

    loadOptions()

    return () => {
      mounted = false
    }
  }, [config])

  useEffect(() => {
    rowsCountRef.current = rows.length
  }, [rows])

  // Загрузка данных
  useEffect(() => {
    let mounted = true

    async function loadRows() {
      const initialLoad = rowsCountRef.current === 0
      setLoading(initialLoad)
      setIsFetching(!initialLoad)
      setError('')

      try {
        const normalizedFilters = normalizeFilters(config.filters, debouncedFilters)
        const offset = (page - 1) * pageSize
        
        const result = await apiMap[config.key].search(
          normalizedFilters,
          pageSize,
          offset,
        )

        if (!mounted) {
          return
        }

        setRows(
          result.items.map((item, index) => ({
            ...item,
            _rowKey: config.getRecordKey(item),
            _rowNumber: offset + index + 1,
          })),
        )
        setTotal(result.count)
      } catch (requestError) {
        if (mounted) {
          setError(
            getApiErrorMessage(requestError, `Не удалось загрузить раздел "${config.title}"`),
          )
        }
      } finally {
        if (mounted) {
          setLoading(false)
          setIsFetching(false)
        }
      }
    }

    loadRows()

    return () => {
      mounted = false
    }
  }, [config, debouncedFilters, page, pageSize])

  const refreshCurrentPage = async () => {
    const initialLoad = rowsCountRef.current === 0
    setLoading(initialLoad)
    setIsFetching(!initialLoad)
    setError('')

    try {
      const result = await apiMap[config.key].search(
        normalizeFilters(config.filters, debouncedFilters),
        pageSize,
        (page - 1) * pageSize,
      )

      setRows(
        result.items.map((item, index) => ({
          ...item,
          _rowKey: config.getRecordKey(item),
          _rowNumber: (page - 1) * pageSize + index + 1,
        })),
      )
      setTotal(result.count)
    } catch (requestError) {
      setError(getApiErrorMessage(requestError, `Не удалось обновить раздел "${config.title}"`))
    } finally {
      setLoading(false)
      setIsFetching(false)
    }
  }

  const handleFormSubmit = async (payload) => {
    setSubmitting(true)
    setError('')
    setSuccess('')

    try {
      if (modalState.mode === 'create') {
        await apiMap[config.key].create(payload)
        setSuccess(`Запись в разделе "${config.title}" успешно создана`)
      } else {
        await apiMap[config.key].update(config.getRecordId(modalState.record), payload)
        setSuccess(`Запись в разделе "${config.title}" успешно обновлена`)
      }

      setModalState({
        show: false,
        mode: 'create',
        record: null,
      })
      await refreshCurrentPage()
    } catch (requestError) {
      setError(getApiErrorMessage(requestError, 'Не удалось сохранить данные'))
    } finally {
      setSubmitting(false)
    }
  }

  const handleDelete = async () => {
    if (!deleteRecord) {
      return
    }

    setSubmitting(true)
    setError('')
    setSuccess('')

    try {
      await apiMap[config.key].remove(config.getRecordId(deleteRecord))
      setDeleteRecord(null)
      setModalState({
        show: false,
        mode: 'create',
        record: null,
      })
      setSuccess(`Запись "${config.getDeleteLabel(deleteRecord)}" удалена`)
      await refreshCurrentPage()
    } catch (requestError) {
      setError(getApiErrorMessage(requestError, 'Не удалось удалить запись'))
    } finally {
      setSubmitting(false)
    }
  }

  // Отладочный вывод
  console.log('Current page:', page, 'Filters:', filters, 'Debounced:', debouncedFilters)

  return (
    <>
      <PageHeader
        eyebrow={config.eyebrow}
        title={config.title}
        description={config.description}
        actions={
          config.createEnabled ? (
            <Button
              onClick={() =>
                setModalState({
                  show: true,
                  mode: 'create',
                  record: null,
                })
              }
            >
              <IoMdAdd className='me-2'/>Создать запись
            </Button>
          ) : null
        }
      />

      {error ? <Alert variant="danger" dismissible>{error}</Alert> : null}
      {success ? <Alert variant="success" dismissible>{success}</Alert> : null}

      <SearchPanel
        fields={config.filters}
        value={filters}
        optionsMap={optionsMap}
        onChange={(name, value) => {
          setFilters((previous) => ({
            ...previous,
            [name]: value,
          }))

          if (page !== 1) {
            setPage(1)
          }
        }}
        onReset={() => {
          const emptyFilters = buildInitialFilters(config.filters, new URLSearchParams())
          setFilters(emptyFilters)
          setDebouncedFilters(emptyFilters)
          setPage(1)
        }}
      />

      <DataTable
        title={config.title}
        subtitle="Открой запись кликом по строке"
        columns={config.columns}
        rows={rows}
        optionsMap={optionsMap}
        page={page}
        pageSize={pageSize}
        total={total}
        loading={loading}
        isFetching={isFetching}
        onPageChange={setPage}
        onPageSizeChange={(nextSize) => {
          setPageSize(nextSize)
          setPage(1)
        }}
        onViewJson={(title, value) =>
          setJsonState({
            show: true,
            title,
            value,
          })
        }
        onRowClick={(record) =>
          setModalState({
            show: true,
            mode: 'edit',
            record,
          })
        }
      />

      <EntityFormModal
        show={modalState.show}
        mode={modalState.mode}
        config={config}
        record={modalState.record}
        optionsMap={optionsMap}
        busy={submitting}
        onHide={() =>
          setModalState({
            show: false,
            mode: 'create',
            record: null,
          })
        }
        onSubmit={handleFormSubmit}
        onDelete={modalState.record && config.deleteEnabled ? () => setDeleteRecord(modalState.record) : null}
      />

      <ConfirmModal
        show={Boolean(deleteRecord)}
        busy={submitting}
        title="Подтвердите удаление"
        body={
          deleteRecord
            ? `Удалить запись "${config.getDeleteLabel(deleteRecord)}"? Это действие нельзя отменить.`
            : ''
        }
        onHide={() => setDeleteRecord(null)}
        onConfirm={handleDelete}
      />

      <JsonViewerModal
        show={jsonState.show}
        title={jsonState.title}
        value={jsonState.value}
        onHide={() =>
          setJsonState({
            show: false,
            title: '',
            value: null,
          })
        }
      />
    </>
  )
}