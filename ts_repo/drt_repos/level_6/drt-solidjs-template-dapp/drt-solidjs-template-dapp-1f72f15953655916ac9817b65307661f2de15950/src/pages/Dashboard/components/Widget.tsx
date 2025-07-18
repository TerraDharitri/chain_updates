import { Card } from 'components/Card';
import { WidgetType } from 'types/widget.types';

export const Widget = ({
  title,
  description,
  reference,
  anchor,
  widget: DrtWidget,
  props = {}
}: WidgetType) => {
  return (
    <Card
      title={title}
      description={description}
      reference={reference}
      anchor={anchor}
    >
      <DrtWidget {...props} />
    </Card>
  );
};
