import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Alert } from './Alert';

describe('Alert Component', () => {
  describe('rendering', () => {
    it('should render error alert with message', () => {
      render(
        <Alert
          type="error"
          message="An error occurred"
        />
      );

      expect(screen.getByText('An error occurred')).toBeInTheDocument();
    });

    it('should render success alert with message', () => {
      render(
        <Alert
          type="success"
          message="Operation completed successfully"
        />
      );

      expect(screen.getByText('Operation completed successfully')).toBeInTheDocument();
    });

    it('should render warning alert with message', () => {
      render(
        <Alert
          type="warning"
          message="This is a warning"
        />
      );

      expect(screen.getByText('This is a warning')).toBeInTheDocument();
    });

    it('should render info alert with message', () => {
      render(
        <Alert
          type="info"
          message="This is informational"
        />
      );

      expect(screen.getByText('This is informational')).toBeInTheDocument();
    });

    it('should render with title and message', () => {
      render(
        <Alert
          type="success"
          title="Success"
          message="Your changes have been saved"
        />
      );

      expect(screen.getByText('Success')).toBeInTheDocument();
      expect(screen.getByText('Your changes have been saved')).toBeInTheDocument();
    });

    it('should render without title when not provided', () => {
      const { container } = render(
        <Alert
          type="info"
          message="Just a message"
        />
      );

      expect(screen.getByText('Just a message')).toBeInTheDocument();
      expect(container.querySelector('h3')).not.toBeInTheDocument();
    });
  });

  describe('styling by type', () => {
    it('should apply correct styling for error type', () => {
      const { container } = render(
        <Alert
          type="error"
          message="Error message"
        />
      );

      const alertDiv = container.firstChild as HTMLElement;
      expect(alertDiv.className).toContain('bg-red-50');
    });

    it('should apply correct styling for success type', () => {
      const { container } = render(
        <Alert
          type="success"
          message="Success message"
        />
      );

      const alertDiv = container.firstChild as HTMLElement;
      expect(alertDiv.className).toContain('bg-green-50');
    });

    it('should apply correct styling for warning type', () => {
      const { container } = render(
        <Alert
          type="warning"
          message="Warning message"
        />
      );

      const alertDiv = container.firstChild as HTMLElement;
      expect(alertDiv.className).toContain('bg-yellow-50');
    });

    it('should apply correct styling for info type', () => {
      const { container } = render(
        <Alert
          type="info"
          message="Info message"
        />
      );

      const alertDiv = container.firstChild as HTMLElement;
      expect(alertDiv.className).toContain('bg-blue-50');
    });
  });

  describe('close button', () => {
    it('should not render close button when onClose is not provided', () => {
      const { container } = render(
        <Alert
          type="info"
          message="Message without close button"
        />
      );

      const buttons = container.querySelectorAll('button');
      expect(buttons).toHaveLength(0);
    });

    it('should render close button when onClose is provided', () => {
      const handleClose = vi.fn();
      const { container } = render(
        <Alert
          type="info"
          message="Message with close button"
          onClose={handleClose}
        />
      );

      const buttons = container.querySelectorAll('button');
      expect(buttons).toHaveLength(1);
    });

    it('should call onClose when close button is clicked', async () => {
      const handleClose = vi.fn();
      const user = userEvent.setup();

      const { container } = render(
        <Alert
          type="error"
          message="Click to close"
          onClose={handleClose}
        />
      );

      const closeButton = container.querySelector('button') as HTMLButtonElement;
      await user.click(closeButton);

      expect(handleClose).toHaveBeenCalledOnce();
    });

    it('should have accessible close button', () => {
      const handleClose = vi.fn();
      const { container } = render(
        <Alert
          type="warning"
          message="Warning"
          onClose={handleClose}
        />
      );

      const closeButton = container.querySelector('button');
      expect(closeButton).toBeInTheDocument();
      expect(closeButton?.type).toBe('button');
    });
  });

  describe('icon rendering', () => {
    it('should render error icon for error type', () => {
      const { container } = render(
        <Alert
          type="error"
          message="Error"
        />
      );

      const svgs = container.querySelectorAll('svg');
      expect(svgs.length).toBeGreaterThan(0);
    });

    it('should render success icon for success type', () => {
      const { container } = render(
        <Alert
          type="success"
          message="Success"
        />
      );

      const svgs = container.querySelectorAll('svg');
      expect(svgs.length).toBeGreaterThan(0);
    });

    it('should render warning icon for warning type', () => {
      const { container } = render(
        <Alert
          type="warning"
          message="Warning"
        />
      );

      const svgs = container.querySelectorAll('svg');
      expect(svgs.length).toBeGreaterThan(0);
    });

    it('should render info icon for info type', () => {
      const { container } = render(
        <Alert
          type="info"
          message="Info"
        />
      );

      const svgs = container.querySelectorAll('svg');
      expect(svgs.length).toBeGreaterThan(0);
    });

    it('should have correct styling for icons', () => {
      const { container } = render(
        <Alert
          type="error"
          message="Error"
        />
      );

      const iconContainer = container.querySelector('div.text-red-400');
      expect(iconContainer).toBeInTheDocument();
      expect(iconContainer?.querySelector('svg')).toBeInTheDocument();
    });
  });

  describe('content variations', () => {
    it('should handle long messages', () => {
      const longMessage =
        'This is a very long error message that should wrap properly without breaking the layout of the alert component.';

      render(
        <Alert
          type="error"
          message={longMessage}
        />
      );

      expect(screen.getByText(longMessage)).toBeInTheDocument();
    });

    it('should handle empty title and message', () => {
      const { container } = render(
        <Alert
          type="info"
          message=""
        />
      );

      expect(container).toBeInTheDocument();
    });

    it('should handle message with special characters', () => {
      const message = 'Error: Invalid input! Check & fix your data.';

      render(
        <Alert
          type="error"
          message={message}
        />
      );

      expect(screen.getByText(message)).toBeInTheDocument();
    });

    it('should handle title with special characters', () => {
      const title = 'Warning: Check here!';
      const message = 'Please review the warning above.';

      render(
        <Alert
          type="warning"
          title={title}
          message={message}
        />
      );

      expect(screen.getByText(title)).toBeInTheDocument();
    });
  });

  describe('accessibility', () => {
    it('should have proper contrast for error type', () => {
      const { container } = render(
        <Alert
          type="error"
          message="Error message"
        />
      );

      const alertDiv = container.firstChild as HTMLElement;
      expect(alertDiv.className).toContain('text-red-800');
    });

    it('should have proper contrast for success type', () => {
      const { container } = render(
        <Alert
          type="success"
          message="Success message"
        />
      );

      const alertDiv = container.firstChild as HTMLElement;
      expect(alertDiv.className).toContain('text-green-800');
    });

    it('should have clickable close button with proper styling', () => {
      const handleClose = vi.fn();
      const { container } = render(
        <Alert
          type="info"
          message="Info"
          onClose={handleClose}
        />
      );

      const closeButton = container.querySelector('button');
      expect(closeButton?.className).toContain('hover:');
    });

    it('should render icons with proper size and color', () => {
      const { container } = render(
        <Alert
          type="success"
          message="Success"
        />
      );

      const icons = container.querySelectorAll('svg');
      icons.forEach((icon) => {
        expect(icon.className.baseVal).toContain('h-5');
        expect(icon.className.baseVal).toContain('w-5');
      });
    });
  });

  describe('multiple alerts', () => {
    it('should render multiple alerts independently', () => {
      const { container } = render(
        <>
          <Alert type="error" message="Error" />
          <Alert type="success" message="Success" />
          <Alert type="warning" message="Warning" />
        </>
      );

      expect(screen.getByText('Error')).toBeInTheDocument();
      expect(screen.getByText('Success')).toBeInTheDocument();
      expect(screen.getByText('Warning')).toBeInTheDocument();

      const alertDivs = container.querySelectorAll('div.bg-red-50, div.bg-green-50, div.bg-yellow-50');
      expect(alertDivs.length).toBeGreaterThanOrEqual(3);
    });
  });
});
