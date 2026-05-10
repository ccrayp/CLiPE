import { Button, Modal } from 'react-bootstrap'
import { FaArrowRightFromBracket } from "react-icons/fa6";
import { FaCheck } from "react-icons/fa";


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
    <Modal show={show} onHide={onHide} centered size="sm">
      <Modal.Header closeButton>
        <Modal.Title>{title}</Modal.Title>
      </Modal.Header>
      <Modal.Body>{body}</Modal.Body>
      <Modal.Footer>
        <Button variant="outline-secondary" onClick={onHide} disabled={busy}>
          <FaArrowRightFromBracket className='me-2'/>Отмена
        </Button>
        <Button variant="danger" onClick={onConfirm} disabled={busy}>
          <FaCheck className='me-2'/>{confirmLabel}
        </Button>
      </Modal.Footer>
    </Modal>
  )
}
