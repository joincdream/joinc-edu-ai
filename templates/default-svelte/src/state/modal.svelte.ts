let isOpen = $state<boolean>(false);
let svgContent = $state<string>('');

export const modalState = {
  get isOpen() {
    return isOpen;
  },
  get svgContent() {
    return svgContent;
  },
  open(svg: string) {
    svgContent = svg;
    isOpen = true;
  },
  close() {
    isOpen = false;
    svgContent = '';
  }
};
