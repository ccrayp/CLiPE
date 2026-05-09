import { Badge, Button, Card, Pagination, Table } from 'react-bootstrap'
import { Link } from 'react-router-dom'
import {
  formatDateTime,
  formatJsonPreview,
  resolveRelationLabel,
} from '../../utils/formatters'

function renderBadge(column, value) {
  return (
    <Badge bg={value ? column.trueVariant ?? 'success' : column.falseVariant ?? 'secondary'}>
      {value ? column.trueLabel ?? 'Да' : column.falseLabel ?? 'Нет'}
    </Badge>
  )
}

function renderLinkedValue(label, to) {
  if (!to || label === '—') {
    return label
  }

  return (
    <Link to={to} onClick={(event) => event.stopPropagation()} className="table-link">
      {label}
    </Link>
  )
}

function buildPaginationItems(page, totalPages, onPageChange) {
  const pages = []
  const start = Math.max(1, page - 2)
  const end = Math.min(totalPages, page + 2)

  if (start > 1) {
    pages.push(
      <Pagination.Item key={1} active={page === 1} onClick={() => onPageChange(1)}>
        1
      </Pagination.Item>,
    )
  }

  if (start > 2) {
    pages.push(<Pagination.Ellipsis key="start-ellipsis" disabled />)
  }

  for (let current = start; current <= end; current += 1) {
    pages.push(
      <Pagination.Item
        key={current}
        active={current === page}
        onClick={() => onPageChange(current)}
      >
        {current}
      </Pagination.Item>,
    )
  }

  if (end < totalPages - 1) {
    pages.push(<Pagination.Ellipsis key="end-ellipsis" disabled />)
  }

  if (end < totalPages) {
    pages.push(
      <Pagination.Item
        key={totalPages}
        active={page === totalPages}
        onClick={() => onPageChange(totalPages)}
      >
        {totalPages}
      </Pagination.Item>,
    )
  }

  return pages
}

function getColumnIdentity(column, index) {
  return column.id ?? `${column.key}-${index}`
}

export function DataTable({
  title,
  subtitle,
  columns,
  rows,
  optionsMap,
  page,
  pageSize,
  total,
  loading,
  onPageChange,
  onViewJson,
  onRowClick,
  isFetching = false,
  onPageSizeChange,
}) {
  const totalPages = Math.max(1, Math.ceil((total || 0) / pageSize))

  return (
    <Card className="surface-card">
      <Card.Header>
        <div className="fw-semibold">{title}</div>
        {subtitle ? <div className="small text-muted">{subtitle}</div> : null}
      </Card.Header>
      <Card.Body className="p-0">
        <div className="table-responsive">
          <Table hover className="mb-0 entity-table align-middle">
            <thead>
              <tr>
                {columns.map((column, index) => (
                  <th key={getColumnIdentity(column, index)}>{column.label}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {loading && rows.length === 0 ? (
                <tr>
                  <td colSpan={columns.length} className="text-center py-5 text-muted">
                    Загрузка данных...
                  </td>
                </tr>
              ) : rows.length === 0 ? (
                <tr>
                  <td colSpan={columns.length} className="text-center py-5 text-muted">
                    По текущим фильтрам записи не найдены
                  </td>
                </tr>
              ) : (
                rows.map((row, index) => (
                  <tr
                    key={row._rowKey ?? index}
                    className="clickable-row"
                    onClick={() => onRowClick(row)}
                  >
                    {columns.map((column, columnIndex) => {
                      const value = row[column.key]
                      const cellKey = getColumnIdentity(column, columnIndex)

                      if (column.type === 'boolean') {
                        return <td key={cellKey}>{renderBadge(column, value)}</td>
                      }

                      if (column.type === 'date') {
                        return <td key={cellKey}>{formatDateTime(value)}</td>
                      }

                      if (column.type === 'relation') {
                        const label = resolveRelationLabel({
                          value,
                          source: column.source,
                          optionValueKey: column.optionValueKey,
                          getOptionLabel: column.getOptionLabel,
                          optionsMap,
                        })
                        const to = column.linkTo ? column.linkTo(value, row) : ''

                        return <td key={cellKey}>{renderLinkedValue(label, to)}</td>
                      }

                      if (column.type === 'link') {
                        return (
                          <td key={cellKey}>
                            {renderLinkedValue(
                              column.getLabel(row),
                              column.getLink(row),
                            )}
                          </td>
                        )
                      }

                      if (column.type === 'json') {
                        return (
                          <td key={cellKey}>
                            <div className="d-flex align-items-center gap-2">
                              <span className="json-preview mono-text">
                                {formatJsonPreview(value)}
                              </span>
                            </div>
                          </td>
                        )
                      }

                      return <td key={cellKey}>{value ?? '—'}</td>
                    })}
                  </tr>
                ))
              )}
            </tbody>
          </Table>
        </div>
      </Card.Body>
      <Card.Footer className="d-flex flex-wrap align-items-center gap-3 table-footer">
        <div className="small text-muted footer-side d-flex align-items-center gap-2 flex-wrap">
          Всего записей: <strong>{total}</strong>
          {onPageSizeChange ? (
            <>
              <span className="ms-2">Строк на странице:</span>
              <select
                className="form-select form-select-sm table-limit-select"
                value={pageSize}
                onChange={(event) => onPageSizeChange(Number(event.target.value))}
              >
                {[10, 20, 50, 100].map((size) => (
                  <option key={size} value={size}>
                    {size}
                  </option>
                ))}
              </select>
            </>
          ) : null}
        </div>
        <Pagination className="mb-0 justify-content-center footer-pagination">
          <Pagination.First
            disabled={page <= 1}
            onClick={() => onPageChange(1)}
          />
          <Pagination.Prev
            disabled={page <= 1}
            onClick={() => onPageChange(Math.max(1, page - 1))}
          />
          {buildPaginationItems(page, totalPages, onPageChange)}
          <Pagination.Next
            disabled={page >= totalPages}
            onClick={() => onPageChange(Math.min(totalPages, page + 1))}
          />
          <Pagination.Last
            disabled={page >= totalPages}
            onClick={() => onPageChange(totalPages)}
          />
        </Pagination>
        <div className="footer-side" />
      </Card.Footer>
    </Card>
  )
}
