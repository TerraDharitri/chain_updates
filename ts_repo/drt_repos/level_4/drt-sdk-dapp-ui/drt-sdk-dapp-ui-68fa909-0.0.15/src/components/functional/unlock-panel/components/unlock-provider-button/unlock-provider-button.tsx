import { Component, h, Prop } from '@stencil/core';
import classNames from 'classnames';
import type { IProviderBase } from 'types/provider.types';
import type { ProviderTypeEnum } from 'types/provider.types';

import { getProviderButtonIcon } from '../../helpers';

const unlockProviderButtonClasses: Record<string, string> = {
  button: 'drt:w-full',
};

@Component({
  tag: 'drt-unlock-provider-button',
  shadow: true,
})
export class UnlockProviderButton {
  @Prop() provider: IProviderBase<ProviderTypeEnum>;
  @Prop() class?: string;

  render() {
    if (!this.provider) {
      return null;
    }

    const icon: HTMLElement | null = !this.provider.iconUrl ? getProviderButtonIcon(this.provider.type) : null;

    return (
      <drt-unlock-button
        iconUrl={this.provider.iconUrl}
        label={this.provider.name}
        type={this.provider.type}
        class={classNames(this.class, unlockProviderButtonClasses.button)}
        icon={icon}
        dataTestId={this.provider.type.toString()}
      />
    );
  }
}
