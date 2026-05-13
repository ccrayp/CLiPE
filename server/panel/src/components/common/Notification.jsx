import { Toast, ToastContainer } from 'react-bootstrap';
import { useState, useEffect } from 'react';

export default function Notification({
  children,
  variant = 'danger',
  dismissible = false,
  className = '',
  delay = 5000,
  autohide = true,
}) {
  const [show, setShow] = useState(false);

  useEffect(() => {
    if (!children) return;
    const timer = setTimeout(() => {
      setShow(true);
    }, 10);

    return () => clearTimeout(timer);
  }, [children]);

  const bgMap = {
    danger: 'danger',
    success: 'success',
    warning: 'warning',
    info: 'info',
    primary: 'primary',
    secondary: 'secondary',
  };

  if (!children) {
    return null;
  }

  return (
    <ToastContainer
      position="bottom-end"
      className={`p-3 ${className}`}
      style={{
        position: 'fixed',
        zIndex: 9999,
        bottom: 0,
        right: 0,
      }}
    >
      <Toast
        show={show}
        animation={true}
        onClose={() => setShow(false)}
        delay={delay}
        autohide={autohide}
        bg={bgMap[variant] || 'danger'}
      >
        <Toast.Header closeButton={dismissible}>
          <strong className="me-auto">
            {variant === 'success'
              ? 'Успешно'
              : variant === 'warning'
              ? 'Предупреждение'
              : variant === 'info'
              ? 'Информация'
              : 'Ошибка'}
          </strong>
        </Toast.Header>

        <Toast.Body
          className={
            ['danger', 'success', 'primary', 'secondary'].includes(
              bgMap[variant]
            )
              ? 'text-white'
              : ''
          }
        >
          {children}
        </Toast.Body>
      </Toast>
    </ToastContainer>
  );
}