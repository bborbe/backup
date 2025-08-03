import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import ActionButtonComponent from '../components/ActionButtonComponent.vue';

describe('ActionButtonComponent', () => {
  it('renders with required props', () => {
    const wrapper = mount(ActionButtonComponent, {
      props: {
        label: 'Test Button',
      },
    });

    expect(wrapper.text()).toBe('Test Button');
    expect(wrapper.find('button').exists()).toBe(true);
  });

  it('applies correct default classes', () => {
    const wrapper = mount(ActionButtonComponent, {
      props: {
        label: 'Test Button',
      },
    });

    const button = wrapper.find('button');
    expect(button.classes()).toContain('btn');
    expect(button.classes()).toContain('btn-primary');
    expect(button.classes()).toContain('btn-normal');
  });

  it('applies variant class correctly', () => {
    const wrapper = mount(ActionButtonComponent, {
      props: {
        label: 'Test Button',
        variant: 'secondary',
      },
    });

    const button = wrapper.find('button');
    expect(button.classes()).toContain('btn-secondary');
  });

  it('applies size class correctly', () => {
    const wrapper = mount(ActionButtonComponent, {
      props: {
        label: 'Test Button',
        size: 'small',
      },
    });

    const button = wrapper.find('button');
    expect(button.classes()).toContain('btn-small');
  });

  it('applies disabled state correctly', () => {
    const wrapper = mount(ActionButtonComponent, {
      props: {
        label: 'Test Button',
        disabled: true,
      },
    });

    const button = wrapper.find('button');
    expect(button.classes()).toContain('btn-disabled');
    expect(button.attributes('disabled')).toBeDefined();
  });

  it('emits click event when clicked and not disabled', async () => {
    const wrapper = mount(ActionButtonComponent, {
      props: {
        label: 'Test Button',
      },
    });

    await wrapper.find('button').trigger('click');

    expect(wrapper.emitted()).toHaveProperty('click');
    expect(wrapper.emitted('click')).toHaveLength(1);
  });

  it('does not emit click event when clicked and disabled', async () => {
    const wrapper = mount(ActionButtonComponent, {
      props: {
        label: 'Test Button',
        disabled: true,
      },
    });

    await wrapper.find('button').trigger('click');

    expect(wrapper.emitted('click')).toBeFalsy();
  });

  it('validates variant prop correctly', () => {
    // Valid variants should not throw
    expect(() => {
      mount(ActionButtonComponent, {
        props: {
          label: 'Test',
          variant: 'primary',
        },
      });
    }).not.toThrow();

    expect(() => {
      mount(ActionButtonComponent, {
        props: {
          label: 'Test',
          variant: 'secondary',
        },
      });
    }).not.toThrow();

    expect(() => {
      mount(ActionButtonComponent, {
        props: {
          label: 'Test',
          variant: 'danger',
        },
      });
    }).not.toThrow();
  });

  it('validates size prop correctly', () => {
    // Valid sizes should not throw
    expect(() => {
      mount(ActionButtonComponent, {
        props: {
          label: 'Test',
          size: 'small',
        },
      });
    }).not.toThrow();

    expect(() => {
      mount(ActionButtonComponent, {
        props: {
          label: 'Test',
          size: 'normal',
        },
      });
    }).not.toThrow();
  });
});