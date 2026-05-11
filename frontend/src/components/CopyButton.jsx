import { useState } from 'react'
import { FiCopy, FiCheck } from 'react-icons/fi'
import './CopyButton.css'

export default function CopyButton({ text, small }) {
  const [copied, setCopied] = useState(false)

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch (err) {
      // фолбэк для небезопасного контекста
      const textarea = document.createElement('textarea')
      textarea.value = text
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }

  return (
    <button
      className={`copy-btn ${small ? 'copy-btn-sm' : ''}`}
      onClick={handleCopy}
      title="Скопировать ID"
    >
      {copied ? <FiCheck /> : <FiCopy />}
    </button>
  )
}