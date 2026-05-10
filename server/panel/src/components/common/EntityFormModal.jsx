/* eslint-disable react-hooks/set-state-in-effect */
import { useEffect, useMemo, useState } from 'react'
import { Alert, Button, Form, Modal, Stack } from 'react-bootstrap'
import { JsonCodeEditor } from './JsonCodeEditor'
import { JsonViewerModal } from './JsonViewerModal'
import { RuleHelpModal } from './RuleHelpModal'
import { MdDeleteOutline } from "react-icons/md";
import { FaArrowRightFromBracket } from "react-icons/fa6";
import { FiSave } from "react-icons/fi";

function toFieldValue(field, record) {
  const rawValue = record?.[field.name]

  if (field.type === 'json') {
    return rawValue ? JSON.stringify(rawValue, null, 2) : ''
  }

  if (field.type === 'boolean') {
    return Boolean(rawValue)
  }

  if (rawValue === null || rawValue === undefined) {
    return ''
  }

  return String(rawValue)
}

function buildInitialState(fields, record) {
  return fields.reduce((accumulator, field) => {
    accumulator[field.name] = toFieldValue(field, record)
    return accumulator
  }, {})
}

function renderOptions(field, optionsMap) {
  const options = optionsMap[field.source] ?? []
  const getLabel =
    field.getOptionLabel ??
    ((item) => item[field.optionLabelKey ?? field.optionValueKey ?? field.name])

  return options.map((item) => {
    const value = item[field.optionValueKey]

    return (
      <option key={`${field.source}-${value}`} value={value}>
        {getLabel(item)}
      </option>
    )
  })
}

function parseValue(field, value) {
  if (field.type === 'boolean') {
    return Boolean(value)
  }

  if (field.type === 'json') {
    if (!value) {
      if (field.required) {
        throw new Error(`Поле "${field.label}" обязательно`)
      }

      return null
    }

    try {
      return JSON.parse(value)
    } catch (error) {
      throw new Error(error.message || `Поле "${field.label}" содержит некорректный JSON`)
    }
  }

  if (field.type === 'number' || field.type === 'select') {
    if (value === '') {
      if (field.required && !field.nullable) {
        throw new Error(`Поле "${field.label}" обязательно`)
      }

      return field.nullable ? null : undefined
    }

    return Number(value)
  }

  if (field.required && !value.trim()) {
    throw new Error(`Поле "${field.label}" обязательно`)
  }

  return value.trim()
}

export function EntityFormModal({
  show,
  mode,
  config,
  record,
  optionsMap,
  busy,
  onHide,
  onSubmit,
  onDelete,
}) {
  const [formData, setFormData] = useState(() =>
    buildInitialState(config.formFields, record),
  )
  const [formError, setFormError] = useState('')
  const [jsonPreview, setJsonPreview] = useState({
    show: false,
    title: '',
    value: null,
  })
  const [showRuleHelp, setShowRuleHelp] = useState(false)

  const title = useMemo(() => {
    if (mode === 'create') {
      return `Создать запись: ${config.title}`
    }

    return `Карточка: ${config.title}`
  }, [config.title, mode])

  useEffect(() => {
    setFormData(buildInitialState(config.formFields, record))
    setFormError('')
  }, [config.formFields, record, show])

  const handleFieldChange = (name, value) => {
    setFormData((previous) => ({
      ...previous,
      [name]: value,
    }))
  }

  const handleSubmit = (event) => {
    event.preventDefault()
    setFormError('')

    try {
      const payload = config.formFields.reduce((accumulator, field) => {
        const parsed = parseValue(field, formData[field.name] ?? '')

        if (parsed !== undefined) {
          accumulator[field.name] = parsed
        }

        return accumulator
      }, {})

      onSubmit(payload)
    } catch (error) {
      setFormError(error.message)
    }
  }

  return (
    <>
      <Modal show={show} onHide={onHide} size="lg" centered>
        <Modal.Header closeButton>
          <Modal.Title>{title}</Modal.Title>
        </Modal.Header>
        <Form onSubmit={handleSubmit}>
          <Modal.Body>
            {formError ? <Alert variant="danger" dismissible>{formError}</Alert> : null}
            <Stack gap={3}>
              {config.formFields.map((field) => (
                <div key={field.name}>
                  <Form.Label className="fw-semibold">{field.label}</Form.Label>
                  {field.type === 'boolean' ? (
                    <Form.Check
                      type="switch"
                      checked={Boolean(formData[field.name])}
                      onChange={(event) =>
                        handleFieldChange(field.name, event.target.checked)
                      }
                      label={field.booleanLabel ?? field.label}
                    />
                  ) : field.type === 'select' ? (
                    <Form.Select
                      value={formData[field.name] ?? ''}
                      onChange={(event) =>
                        handleFieldChange(field.name, event.target.value)
                      }
                    >
                      <option value="">
                        {field.nullable ? 'Не выбрано' : 'Выберите значение'}
                      </option>
                      {renderOptions(field, optionsMap)}
                    </Form.Select>
                  ) : field.type === 'json' ? (
                    <div>
                      {field.showRuleHelp ? (
                        <div className="mb-2">
                          <Button
                            type="button"
                            variant="outline-secondary"
                            size="sm"
                            onClick={() => setShowRuleHelp(true)}
                          >
                            Как описывать условия?
                          </Button>
                        </div>
                      ) : null}
                      <JsonCodeEditor
                        value={formData[field.name] ?? ''}
                        onChange={(nextValue) => handleFieldChange(field.name, nextValue)}
                      />
                    </div>
                  ) : (
                    <Form.Control
                      type={field.type === 'number' ? 'number' : 'text'}
                      value={formData[field.name] ?? ''}
                      onChange={(event) =>
                        handleFieldChange(field.name, event.target.value)
                      }
                    />
                  )}
                  {field.helpText ? (
                    <Form.Text className="text-muted">{field.helpText}</Form.Text>
                  ) : null}
                </div>
              ))}
            </Stack>
          </Modal.Body>
          <Modal.Footer className="justify-content-between">
            <div>
              {record && onDelete ? (
                <Button variant="outline-danger" onClick={onDelete} disabled={busy}>
                  <MdDeleteOutline className='me-2'/>Удалить
                </Button>
              ) : null}
            </div>
            <div className="d-flex gap-2">
              <Button variant="outline-secondary" onClick={onHide} disabled={busy}>
                <FaArrowRightFromBracket className='me-2'/>Закрыть
              </Button>
              <Button type="submit" disabled={busy}>
                <FiSave className='me-2'/>{mode === 'create' ? 'Создать' : 'Сохранить'}
              </Button>
            </div>
          </Modal.Footer>
        </Form>
      </Modal>

      <JsonViewerModal
        show={jsonPreview.show}
        title={jsonPreview.title}
        value={jsonPreview.value}
        onHide={() =>
          setJsonPreview({
            show: false,
            title: '',
            value: null,
          })
        }
      />

      <RuleHelpModal show={showRuleHelp} onHide={() => setShowRuleHelp(false)} />
    </>
  )
}
