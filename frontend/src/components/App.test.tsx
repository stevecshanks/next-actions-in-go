import React from "react";
// @ts-ignore See: https://github.com/testing-library/react-testing-library/issues/610
import { render, waitFor } from "@testing-library/react";
import fetchMock from "jest-fetch-mock";
import MockDate from "mockdate";
import App from "./App";
import API_SUCCESS_RESPONSE from "../../../contracts/api_success_response.json";
import API_ERROR_RESPONSE from "../../../contracts/api_error_response.json";

const NOW = new Date(2020, 0, 15, 10, 30, 0);
const LAST_YEAR = new Date(2019, 11, 31, 10, 30, 0);

beforeEach(() => MockDate.set(NOW));

afterEach(() => MockDate.reset());

test("renders the actions returned by the API", async () => {
  fetchMock.mockResponse(JSON.stringify(API_SUCCESS_RESPONSE));

  const { findByText, queryByText } = render(<App />);

  const firstAction = await findByText("My First Action");
  expect(firstAction).toBeInTheDocument();
  const secondAction = queryByText("My Second Action");
  expect(secondAction).toBeInTheDocument();
  const todoAction = queryByText("Todo Action");
  expect(todoAction).toBeInTheDocument();
  const projectAction = queryByText("Project Action");
  expect(projectAction).toBeInTheDocument();
});

test("correctly handles due by datetimes from the API", async () => {
  fetchMock.mockResponse(JSON.stringify(API_SUCCESS_RESPONSE));

  const { findByText, queryByText } = render(<App />);

  const firstAction = await findByText("1/1/2020");
  expect(firstAction).toBeInTheDocument();
  const secondAction = queryByText("1/15/2020");
  expect(secondAction).toBeInTheDocument();
  const badDueDate = queryByText("1/1/1970");
  expect(badDueDate).not.toBeInTheDocument();
});

test("includes the count of overdue and due soon actions in the window title", async () => {
  fetchMock.mockResponse(JSON.stringify(API_SUCCESS_RESPONSE));

  render(<App />);

  await waitFor(() => expect(document.title).toEqual("(2) Next Actions"));
});

test("does not include action count in title if no overdue or due soon actions", async () => {
  fetchMock.mockResponse(JSON.stringify(API_SUCCESS_RESPONSE));
  MockDate.set(LAST_YEAR);

  const { findByText } = render(<App />);

  // Wait for actions to actually render, otherwise the test just picks up the initial title
  await findByText("My First Action");
  await waitFor(() => expect(document.title).toEqual("Next Actions"));
});

test("renders an error if the API call fails", async () => {
  fetchMock.mockReject();

  const { findByText } = render(<App />);

  const error = await findByText("An error occurred");
  expect(error).toBeInTheDocument();
});

test("renders errors returned from the API", async () => {
  fetchMock.mockResponse(JSON.stringify(API_ERROR_RESPONSE), { status: 500 });

  const { findByText } = render(<App />);

  const error = await findByText(
    "An error occurred: request to /members/me/cards returned status code 404"
  );
  expect(error).toBeInTheDocument();
});

test("includes an indicator in the window title when errors are returned from the API", async () => {
  fetchMock.mockResponse(JSON.stringify(API_ERROR_RESPONSE), { status: 500 });

  render(<App />);

  await waitFor(() => expect(document.title).toEqual("[ERROR] Next Actions"));
});
