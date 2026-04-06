async function bootstrap() {
  const summary = document.querySelector('#summary');
  const sourceList = document.querySelector('#source-list');
  const sourceTitle = document.querySelector('#source-title');
  const sourceMeta = document.querySelector('#source-meta');
  const contentPanel = document.querySelector('#content-panel');
  const contentState = document.querySelector('#content-state');
  const sourceContent = document.querySelector('#source-content');

  function renderMeta(parts) {
    const items = parts.filter(Boolean).map((value) => {
      const item = document.createElement('span');
      item.textContent = value;
      return item;
    });
    sourceMeta.replaceChildren(...items);
  }

  function renderPanel(state, message, content = '') {
    contentPanel.dataset.state = state;

    if (state === 'ready') {
      contentState.hidden = true;
      sourceContent.hidden = false;
      sourceContent.textContent = content;
      return;
    }

    sourceContent.hidden = true;
    sourceContent.textContent = '';
    contentState.hidden = false;
    contentState.textContent = message;
  }

  const response = await fetch('/api/sources');
  if (!response.ok) {
    throw new Error(`Request failed (${response.status})`);
  }
  const payload = await response.json();
  const sources = payload.sources || [];
  const generatedAt = payload.generatedAt ? new Date(payload.generatedAt).toLocaleString() : 'unknown';

  const total = document.createElement('strong');
  total.textContent = `${payload.totalSources || sources.length} sources`;
  const root = document.createElement('p');
  root.textContent = `root · ${payload.root || 'unknown'}`;
  const generated = document.createElement('p');
  generated.textContent = `updated · ${generatedAt}`;
  summary.replaceChildren(total, root, generated);

  if (sources.length === 0) {
    sourceTitle.textContent = 'No sources found';
    renderMeta(['allowlist is empty']);
    renderPanel('empty', 'Add an allowlisted file such as AGENTS.md, CLAUDE.md, .claude/, .cursor/rules/, or .sisyphus/.');
    return;
  }

  function render(selectedSource) {
    sourceTitle.textContent = selectedSource.label;
    renderMeta([selectedSource.family, selectedSource.relativePath, `${selectedSource.sizeBytes} bytes`]);
    renderPanel('ready', '', selectedSource.content);
  }

  sources.forEach((item, index) => {
    const listItem = document.createElement('li');
    const button = document.createElement('button');
    const label = document.createElement('strong');
    const meta = document.createElement('small');
    button.type = 'button';
    button.className = 'source-button';
    label.textContent = item.label;
    meta.textContent = `${item.family} · ${item.relativePath}`;
    button.append(label, meta);
    button.addEventListener('click', () => {
      document.querySelectorAll('.source-button').forEach((entry) => entry.classList.remove('is-active'));
      button.classList.add('is-active');
      render(item);
    });

    if (index === 0) {
      button.classList.add('is-active');
      render(item);
    }

    listItem.appendChild(button);
    sourceList.appendChild(listItem);
  });
}

bootstrap().catch((error) => {
  const sourceTitle = document.querySelector('#source-title');
  const sourceMeta = document.querySelector('#source-meta');
  const contentPanel = document.querySelector('#content-panel');
  const contentState = document.querySelector('#content-state');
  const sourceContent = document.querySelector('#source-content');
  sourceTitle.textContent = 'Load failed';
  sourceMeta.replaceChildren();
  contentPanel.dataset.state = 'error';
  sourceContent.hidden = true;
  sourceContent.textContent = '';
  contentState.hidden = false;
  contentState.textContent = error.message;
});
