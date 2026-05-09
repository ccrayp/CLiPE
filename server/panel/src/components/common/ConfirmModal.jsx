import { Button, Modal } from 'react-bootstrap'

export function ConfirmModal({
  show,
  title,
  body,
  confirmLabel = 'Удалить',
  onConfirm,
  onHide,
  busy = false,
}) {
  return (
    <Modal show={show} onHide={onHide} centered>
      <Modal.Header closeButton>
        <Modal.Title>{title}</Modal.Title>
      </Modal.Header>
      <Modal.Body>{body}</Modal.Body>
      <Modal.Footer>
        <Button variant="outline-secondary" onClick={onHide} disabled={busy}>
          Отмена
        </Button>
        <Button variant="danger" onClick={onConfirm} disabled={busy}>
          {confirmLabel}
        </Button>
      </Modal.Footer>
    </Modal>
  )
}
