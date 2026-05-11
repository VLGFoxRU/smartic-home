import { FiAlertTriangle } from 'react-icons/fi'
import './ConfirmDialog.css'

export default function ConfirmDialog({
  title = 'Подтверждение',
  message = 'Вы уверены?',
  confirmLabel = 'Удалить',
  onConfirm,
  onCancel,
}) {
  return (
    <div className="confirm-overlay" onClick={onCancel}>
      <div className="confirm-box" onClick={(e) => e.stopPropagation()}>
        <FiAlertTriangle className="confirm-icon" />
        <h3 className="confirm-title">{title}</h3>
        <p className="confirm-message">{message}</p>
        <div className="confirm-buttons">
          <button className="confirm-btn cancel" onClick={onCancel}>
            Отмена
          </button>
          <button className="confirm-btn danger" onClick={onConfirm}>
            {confirmLabel}
          </button>
        </div>
      </div>
    </div>
  )
}