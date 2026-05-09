import editorModule from 'react-simple-code-editor'
import Prism from 'prismjs'
import 'prismjs/components/prism-json'
import 'prismjs/themes/prism.css'

const Editor = editorModule.default ?? editorModule

function getJsonError(value) {
  if (!value.trim()) {
    return ''
  }

  try {
    JSON.parse(value)
    return ''
  } catch (error) {
    return error.message
  }
}

export function JsonCodeEditor({ value, onChange }) {
  const error = getJsonError(value)

  return (
    <div className={`code-editor-shell ${error ? 'has-error' : ''}`}>
      <Editor
        value={value}
        onValueChange={onChange}
        highlight={(code) => Prism.highlight(code, Prism.languages.json, 'json')}
        padding={12}
        textareaId="json-code-editor"
        className="json-code-editor mono-text"
        style={{
          fontFamily: '"SFMono-Regular", "JetBrains Mono", Consolas, monospace',
          fontSize: 14,
          minHeight: 220,
        }}
      />
      {error ? <div className="json-error-text">Ошибка JSON: {error}</div> : null}
    </div>
  )
}
