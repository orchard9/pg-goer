# Claude Development Guide

## Memory Management

### Memory Structure
The project maintains four memory files in the `.memory/` directory:

1. **state.md**: Lists every file in the project with a one-line description of its purpose
2. **working-memory.md**: Tracks work progress
   - Recent: 1-3 lines describing completed work
   - Current: 1-3 lines describing work in progress
   - Future: 1-3 lines describing upcoming work
3. **semantic-memory.md**: Contains simple factual statements about the project
4. **vision.md**: Defines the stable long-term vision for the project

## How to Handle Coding Task Requests

### 0. Gather Requirements
Ask questions until you have 100% confidence in understanding:
- Project specifications
- Requirements
- Expected outcomes
- Edge cases

### 1. Create Project Header
Write a clear header that describes the project's purpose and scope.

### 2. Load Memory and Follow Principles
- First: Load all relevant memory files
- Second: Build elegantly following DRY principles
- Third: Ensure flexibility for future development
- Fourth: Maintain strong adherence to the project vision

### 3. Create 80/20 Tests
When relevant to the task:
- Write simple tests that cover 80% of use cases with 20% effort
- Focus on core functionality first
- Outline expected behavior clearly

### 4. Review Test Coverage
Double-check that tests:
- Follow the 80/20 coverage principle
- Adhere to specifications and vision
- Identify and address any technical debt

### 5. Scaffold Code Structure
Create the file structure with 1-3 line descriptions for each file explaining:
- The file's purpose
- What functionality it provides
- How it fits into the overall architecture

### 6. Validate Architecture
Apply the 80/20 rule to pressure test:
- How code will be called and used
- Whether the architecture supports the requirements
- If changes are needed before implementation

### 7. Implement Code
Write and refine code until:
- All tests pass
- Code meets quality standards
- Implementation matches the design

### 8. User Acceptance Testing
Perform final validation:
- Test from the user's perspective
- Verify all requirements are met
- Ensure the solution is intuitive and reliable

### 9. Update Memory
After completing the task:
- Update state.md with new files
- Update working-memory.md with completed work
- Add new facts to semantic-memory.md if applicable
- Ensure vision.md remains accurate