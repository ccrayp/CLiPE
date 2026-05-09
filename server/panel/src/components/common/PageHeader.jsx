import { Breadcrumb, Col, Row } from 'react-bootstrap'

export function PageHeader({ eyebrow, title, description, actions }) {
  return (
    <section className="page-hero mb-4">
      <Row className="align-items-center g-3">
        <Col lg={8}>
          <Breadcrumb className="mb-1 page-breadcrumbs">
            {eyebrow ? <Breadcrumb.Item active>{eyebrow}</Breadcrumb.Item> : null}
            <Breadcrumb.Item active>{title}</Breadcrumb.Item>
          </Breadcrumb>
          {description ? <p className="mb-0 text-secondary small">{description}</p> : null}
        </Col>
        {actions ? (
          <Col lg={4}>
            <div className="d-flex justify-content-lg-end">{actions}</div>
          </Col>
        ) : null}
      </Row>
    </section>
  )
}
