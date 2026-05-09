import { Modal, Table } from 'react-bootstrap'

const examples = {
  gid: `{
  "type": "gid",
  "operator": "equals",
  "value": 1001
}`,
  groups: `{
  "type": "groups",
  "operator": "contains",
  "value": "sudo"
}`,
  ip: `{
  "type": "ip",
  "operator": "in",
  "value": "192.168.1.0/24"
}`,
  hostname: `{
  "type": "hostname",
  "operator": "regex",
  "value": "^prod-.*"
}`,
  timestamp: `{
  "type": "timestamp",
  "operator": "between",
  "value": "09:00-18:00"
}`,
  weekday: `{
  "type": "weekday",
  "operator": "in",
  "value": ["mon", "tue", "wed"]
}`,
}

export function RuleHelpModal({ show, onHide }) {
  return (
    <Modal show={show} onHide={onHide} size="xl" centered>
      <Modal.Header closeButton>
        <Modal.Title>Как описывать условия в правилах</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <p>
          Правило состоит из массива <code>conditions</code> и поля <code>effect</code>.
          Все условия объединяются через <strong>AND</strong>: если хотя бы одно не
          выполняется, доступ запрещается.
        </p>

        <pre className="help-code-block">
{`{
  "conditions": [
    {
      "type": "ip",
      "operator": "in",
      "value": "192.168.1.0/24"
    }
  ],
  "effect": true
}`}
        </pre>

        <div className="fw-semibold mb-2">Структура одного условия</div>
        <pre className="help-code-block">
{`{
  "type": "<тип>",
  "operator": "<оператор>",
  "value": "<значение>"
}`}
        </pre>

        <Table bordered size="sm" className="mb-4">
          <thead>
            <tr>
              <th>Поле</th>
              <th>Описание</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>type</td>
              <td>Тип проверяемого параметра</td>
            </tr>
            <tr>
              <td>operator</td>
              <td>Оператор сравнения</td>
            </tr>
            <tr>
              <td>value</td>
              <td>Значение для сравнения</td>
            </tr>
          </tbody>
        </Table>

        <div className="fw-semibold mb-2">Поддерживаемые условия</div>

        <div className="mb-3">
          <div className="fw-semibold">gid</div>
          <div className="small text-muted mb-2">Операторы: equals, not_equals</div>
          <pre className="help-code-block">{examples.gid}</pre>
        </div>

        <div className="mb-3">
          <div className="fw-semibold">groups</div>
          <div className="small text-muted mb-2">Операторы: equals, contains</div>
          <pre className="help-code-block">{examples.groups}</pre>
        </div>

        <div className="mb-3">
          <div className="fw-semibold">ip</div>
          <div className="small text-muted mb-2">Операторы: equals, in, not_in</div>
          <pre className="help-code-block">{examples.ip}</pre>
        </div>

        <div className="mb-3">
          <div className="fw-semibold">hostname</div>
          <div className="small text-muted mb-2">Операторы: equals, not_equals, regex</div>
          <pre className="help-code-block">{examples.hostname}</pre>
        </div>

        <div className="mb-3">
          <div className="fw-semibold">timestamp</div>
          <div className="small text-muted mb-2">Оператор: between</div>
          <pre className="help-code-block">{examples.timestamp}</pre>
        </div>

        <div className="mb-0">
          <div className="fw-semibold">weekday</div>
          <div className="small text-muted mb-2">Операторы: equals, in</div>
          <pre className="help-code-block">{examples.weekday}</pre>
        </div>
      </Modal.Body>
    </Modal>
  )
}
