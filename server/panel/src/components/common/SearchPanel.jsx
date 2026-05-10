import { Button, Card, Col, Form, Row } from 'react-bootstrap'
import { TfiReload } from "react-icons/tfi";


function renderOptions(field, optionsMap) {
  const options = optionsMap[field.source] ?? []
  const labelBuilder =
    field.getOptionLabel ??
    ((item) => item[field.optionLabelKey ?? field.optionValueKey ?? field.name])

  return options.map((item) => {
    const value = item[field.optionValueKey ?? field.name]
    return (
      <option key={`${field.source}-${value}`} value={value}>
        {labelBuilder(item)}
      </option>
    )
  })
}

export function SearchPanel({ fields, value, onChange, onReset, optionsMap }) {
  return (
    <Card className="surface-card mb-3">
      <Card.Body>
        <Row className="g-3 align-items-end">
          {fields.map((field) => (
            <Col key={field.name} md={6} xl={4}>
              <Form.Group controlId={`filter-${field.name}`}>
                <Form.Label>{field.label}</Form.Label>
                {field.type === 'triBoolean' ? (
                  <Form.Select
                    value={value[field.name] ?? ''}
                    onChange={(event) => onChange(field.name, event.target.value)}
                  >
                    <option value="">{field.allLabel ?? 'Все'}</option>
                    <option value="true">{field.trueFilterLabel ?? 'Да'}</option>
                    <option value="false">{field.falseFilterLabel ?? 'Нет'}</option>
                  </Form.Select>
                ) : field.type === 'select' ? (
                  <Form.Select
                    value={value[field.name] ?? ''}
                    onChange={(event) => onChange(field.name, event.target.value)}
                  >
                    <option value="">Все</option>
                    {renderOptions(field, optionsMap)}
                  </Form.Select>
                ) : (
                  <Form.Control
                    type={field.type === 'number' ? 'number' : 'text'}
                    value={value[field.name] ?? ''}
                    onChange={(event) => onChange(field.name, event.target.value)}
                    placeholder="Начните ввод..."
                  />
                )}
              </Form.Group>
            </Col>
          ))}
          <Col xs={12} className="d-flex justify-content-end">
            <Button variant="outline-secondary" onClick={onReset}>
              <TfiReload className='me-2'/>Сбросить фильтры
            </Button>
          </Col>
        </Row>
      </Card.Body>
    </Card>
  )
}
