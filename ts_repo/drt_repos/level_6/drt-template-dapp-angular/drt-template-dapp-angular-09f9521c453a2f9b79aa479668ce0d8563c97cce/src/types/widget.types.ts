export type WidgetType<T = any> = {
  anchor?: string;
  description?: string;
  props?: any;
  reference: string;
  title: string;
  widget: any;
}; 