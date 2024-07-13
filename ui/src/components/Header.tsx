import { LabelLarge } from 'baseui/typography';

interface HeaderProps {
  title: string;
}

export function Header(props: HeaderProps) {
  return (
    <LabelLarge>
      {props.title}
    </LabelLarge>
  );
}
