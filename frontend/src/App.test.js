import { render, screen } from '@testing-library/react';
import App from './App';

test('renders system entry title', () => {
  render(<App />);
  const titleElement = screen.getByText(/System Entry/i);
  expect(titleElement).toBeInTheDocument();
});
