export const containerStyle = `
    z-index: 9999;
    position: fixed;
    visibility: hidden;
    right: 0;
    left: auto;
    border: none;
    border-radius: 12px;
    border-bottom-right-radius: 0;
    border-top-right-radius: 0;
    height: calc(100vh - 8px - 64px);
    transition: all 400ms ease;
    width: 420px;
    max-width: calc(100vw - 8px);
    background: black;
    overflow: hidden;
    top: 64px;
  `;

export const headerStyle = `
    color: #fff;
    cursor: pointer;
    height: 80px;
    border-radius: 12px;
    display: flex;
    justify-content: space-between;
    align-items: center;
  `;

export const bodyStyle = `
    height: calc(100vh - 80px - 64px - 8px);
    width: 420px;
    max-width: calc(100vw - 8px);
    overflow: hidden;
  `;

export const iframeStyle = `
    width: 100%;
    height: 100%;
    border: none;
  `;

export const headingElementStyle = `
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 12px;
  width: 100%;
  position: relative;
  padding: 0 24px;
  height: 80px;
`;

export const titleElementStyle = `
  font-size: 20px;
  font-family: 'Roobert Medium', sans-serif;
  color: #E5E5E5;
  transition: all 400ms ease;
  line-height: 1;
  user-select: none;
`;

export const toggleIconElementStyle = `
  margin-left: auto;
  transition: all 400ms ease;
`;
