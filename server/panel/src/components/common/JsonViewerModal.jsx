import { Modal } from 'react-bootstrap'

export function JsonViewerModal({ show, title, value, onHide }) {
  return (
    <Modal show={show} onHide={onHide} size="lg" centered>
      <Modal.Header closeButton>
        <Modal.Title>{title}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <pre className="bg-dark text-light rounded-4 p-4 mono-text overflow-auto">
          {JSON.stringify(value, null, 2)}
        </pre>
      </Modal.Body>
    </Modal>
  )
}
