import { Component, Host, h } from '@stencil/core';

@Component({
  tag: 'tag-143',
  styleUrl: 'tag-143.scss',
  shadow: true,
})
export class Tag143 {
  render() {
    return (
      <Host>
        <slot></slot>
      </Host>
    );
  }
}
